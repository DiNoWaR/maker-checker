package service

import (
	"encoding/json"
	"github.com/dinowar/maker-checker/internal/pkg/domain/model"
	"go.uber.org/zap"
)

type SenderService struct {
	logger *zap.Logger
}

func NewSenderService(logger *zap.Logger) *SenderService {
	return &SenderService{
		logger: logger,
	}
}

// SendMessage for the simplicity we're emulating sending messages by printing them on console
func (service *SenderService) SendMessage(message model.Message) {
	jsonMsg, _ := json.Marshal(message)
	service.logger.Info("message sent", zap.String("message", string(jsonMsg)))
}
