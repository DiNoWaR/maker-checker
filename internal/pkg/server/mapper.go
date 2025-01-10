package server

import (
	"github.com/dinowar/maker-checker/internal/pkg/domain/model"
	"strings"
)

func MapMessageStatusToDb(input string) model.MessageStatus {
	status := strings.ToLower(input)
	if status == "pending" {
		return model.StatusPending
	}
	if status == "approved" {
		return model.StatusApproved
	}
	if status == "rejected" {
		return model.StatusRejected
	}
	return model.StatusAll
}
