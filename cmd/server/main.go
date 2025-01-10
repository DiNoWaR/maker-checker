package main

import (
	"context"
	"fmt"
	"github.com/dinowar/maker-checker/internal/pkg/config"
	"github.com/dinowar/maker-checker/internal/pkg/server"
	"github.com/dinowar/maker-checker/internal/pkg/service"
	"github.com/dinowar/maker-checker/internal/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/sethvargo/go-envconfig"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	ctx := context.Background()
	serviceConfig := &config.ServiceConfig{}
	if configErr := envconfig.Process(ctx, serviceConfig); configErr != nil {
		logger.Fatal("failed to init config", zap.Error(configErr))
	}
	db, dbErr := util.InitDB(
		serviceConfig.DBConfig.Host,
		serviceConfig.DBConfig.Port,
		serviceConfig.DBConfig.Database,
		serviceConfig.DBConfig.Username,
		serviceConfig.DBConfig.Password,
	)
	if dbErr != nil {
		logger.Fatal("failed to init database", zap.Error(dbErr))
	}

	repService := service.NewRepositoryService(db)
	logService := service.NewLogService(logger)
	senderService := service.NewSenderService(logger)
	appServer := server.NewAppServer(repService, senderService, logService, serviceConfig)

	router := gin.Default()

	// Users route
	router.POST("/messages", appServer.CreateMessage)

	// Secured routes for moderators
	protected := router.Group("/messages", server.AuthMiddleware())
	{
		protected.GET("/", appServer.GetMessages)
		protected.PUT("/:id/approve", appServer.ApproveMessage)
		protected.PUT("/:id/reject", appServer.RejectMessage)
	}

	logger.Info(fmt.Sprintf("service starting on port: %s", serviceConfig.ServicePort))
	serverStartErr := router.Run(fmt.Sprintf("%s:%s", serviceConfig.ServiceHost, serviceConfig.ServicePort))
	if serverStartErr != nil {
		logger.Fatal("failed to start server", zap.Error(serverStartErr))
	}
}
