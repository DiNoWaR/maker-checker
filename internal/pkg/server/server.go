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
	sender *service.SenderService
	logger *service.LogService
	config *config.ServiceConfig
}

func NewAppServer(rep *service.RepositoryService, senderService *service.SenderService, logger *service.LogService, config *config.ServiceConfig) *Server {
	return &Server{
		rep:    rep,
		sender: senderService,
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
		return
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
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}

func (server *Server) ApproveMessage(context *gin.Context) {
	messageId := context.Param("id")

	message, messageErr := server.rep.GetMessageById(messageId)
	if messageErr != nil {
		server.logger.LogError("dbErr", messageErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if message.Status != model.StatusPending {
		context.JSON(http.StatusConflict, gin.H{
			"error":   "Message has already been processed",
			"message": "The message is not in a pending state and cannot be updated.",
		})
		return
	}

	dbErr := server.rep.UpdateMessage(messageId, model.StatusApproved)
	if dbErr != nil {
		server.logger.LogError("dbErr", dbErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	// Sending approved message to recipient
	server.sender.SendMessage(*message)

	context.JSON(http.StatusOK, gin.H{
		"messageStatus": "Message approved",
		"messageId":     messageId,
		"content":       message.Content,
	})
}

func (server *Server) RejectMessage(context *gin.Context) {
	messageId := context.Param("id")

	message, messageErr := server.rep.GetMessageById(messageId)
	if messageErr != nil {
		server.logger.LogError("dbErr", messageErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	if message.Status != model.StatusPending {
		context.JSON(http.StatusConflict, gin.H{
			"error":   "Message has already been processed",
			"message": "The message is not in a pending state and cannot be updated.",
		})
		return
	}

	dbErr := server.rep.UpdateMessage(messageId, model.StatusRejected)
	if dbErr != nil {
		server.logger.LogError("dbErr", dbErr)
		context.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"messageStatus": "Message rejected",
		"messageId":     messageId,
		"content":       message.Content,
	})
}
