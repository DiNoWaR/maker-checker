package model

import (
	"time"
)

type Message struct {
	Id          string
	SenderId    string
	RecipientId string
	Content     string
	Status      MessageStatus
	Ts          time.Time
}

type MessageStatus string

const (
	StatusPending  MessageStatus = "PENDING"
	StatusApproved MessageStatus = "APPROVED"
	StatusRejected MessageStatus = "REJECTED"
	StatusAll      MessageStatus = "ALL"
)
