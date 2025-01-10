package server

import (
	"github.com/dinowar/maker-checker/internal/pkg/config"
	"github.com/dinowar/maker-checker/internal/pkg/domain/model"
	"github.com/dinowar/maker-checker/internal/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type CreateMessageRequest struct {
	SenderId    string `json:"senderId" binding:"required"`
	RecipientId string `json:"recipientId" binding:"required"`
	Content     string `json:"content" binding:"required"`
}

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
	var req CreateMessageRequest
	if bindErr := context.ShouldBindJSON(&req); bindErr != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request body: " + bindErr.Error(),
		})
		return
	}

	message := &model.Message{
		Id:          uuid.NewString(),
		SenderId:    req.SenderId,
		RecipientId: req.RecipientId,
		Content:     req.Content,
		Status:      model.StatusPending,
		Ts:          time.Now().In(time.UTC),
	}

	dbErr := server.rep.SaveMessage(message)
	if dbErr != nil {
		server.logger.LogError("dbErr", dbErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
	}

	context.JSON(http.StatusOK, gin.H{
		"id":          message.Id,
		"senderId":    req.SenderId,
		"recipientId": req.RecipientId,
		"content":     req.Content,
		"status":      message.Status,
	})
}

func (server *Server) GetMessages(context *gin.Context) {
	status := context.DefaultQuery("status", "ALL")
	messages, dbErr := server.rep.GetMessages(MapMessageStatusToDb(status))
	if dbErr != nil {
		server.logger.LogError("dbErr", dbErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
	}

	context.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
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
