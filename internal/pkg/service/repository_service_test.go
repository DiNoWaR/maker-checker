package service

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dinowar/maker-checker/internal/pkg/domain/model"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSaveMessage(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rep := NewRepositoryService(db)

	msg := &model.Message{
		Id:          "123",
		SenderId:    "sender1",
		RecipientId: "recipient1",
		Content:     "Hello, world!",
		Status:      model.StatusPending,
		Ts:          time.Now(),
	}

	mock.ExpectExec(`INSERT INTO messages`).
		WithArgs(msg.Id, msg.SenderId, msg.RecipientId, msg.Content, msg.Status, msg.Ts).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := rep.SaveMessage(msg)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}

func TestGetMessageById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rep := NewRepositoryService(db)

	msg := &model.Message{
		Id:          "123",
		SenderId:    "sender1",
		RecipientId: "recipient1",
		Content:     "Hello, world!",
		Status:      model.StatusPending,
		Ts:          time.Date(2025, time.January, 10, 14, 49, 58, 129680000, time.Local),
	}

	rows := sqlmock.NewRows([]string{"id", "sender_id", "recipient_id", "content", "status", "ts"}).
		AddRow(msg.Id, msg.SenderId, msg.RecipientId, msg.Content, msg.Status, msg.Ts)

	mock.ExpectQuery(`SELECT id, sender_id, recipient_id, content, status, ts FROM messages WHERE id = \$1`).
		WithArgs(msg.Id).
		WillReturnRows(rows)

	result, err := rep.GetMessageById("123")
	assert.NoError(t, err)
	assert.Equal(t, msg, result)
	mock.ExpectationsWereMet()
}

func TestUpdateMessage(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	rep := NewRepositoryService(db)

	messageId := "123"
	status := model.StatusApproved

	mock.ExpectExec(`UPDATE messages SET status = \$1 WHERE id = \$2`).
		WithArgs(status, messageId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := rep.UpdateMessage(messageId, status)
	assert.NoError(t, err)
	mock.ExpectationsWereMet()
}
