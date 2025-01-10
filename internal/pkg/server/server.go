package server

import (
	"github.com/dinowar/maker-checker/internal/pkg/config"
	"github.com/dinowar/maker-checker/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	rep    *service.RepositoryService
	logger *service.LogService
	config *config.ServiceConfig
}

func NewAppServer(rep *service.RepositoryService, logger *service.LogService, config *config.ServiceConfig) *Server {
	return &Server{
		rep:    rep,
		logger: logger,
		config: config,
	}
}

func (server *Server) CreateMessage(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "Message created",
	})
}

func (server *Server) GetMessages(context *gin.Context) {
	status := context.DefaultQuery("status", "ALL")

	if status == "PENDING" {
		context.JSON(http.StatusOK, gin.H{
			"messages": []string{"Message 1", "Message 2"},
			"status":   status,
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"messages": []string{"No messages with the specified status"},
			"status":   status,
		})
	}
}

func (server *Server) ApproveMessage(context *gin.Context) {
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"message": "Message approved",
		"id":      id,
	})
}

func (server *Server) RejectMessage(context *gin.Context) {
	id := context.Param("id")
	context.JSON(http.StatusOK, gin.H{
		"message": "Message rejected",
		"id":      id,
	})
}
