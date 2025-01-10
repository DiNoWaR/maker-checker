package service

import (
	"database/sql"
	. "github.com/dinowar/maker-checker/internal/pkg/domain/model"
)

type RepositoryService struct {
	db *sql.DB
}

func NewRepositoryService(db *sql.DB) *RepositoryService {
	return &RepositoryService{db: db}
}

func (rep *RepositoryService) SaveMessage(msg *Message) error {
	if msg.Status == "" {
		msg.Status = StatusPending
	}
	_, saveErr := rep.db.Exec(
		`INSERT INTO messages (id, sender_id, recipient_id, content, status, ts)
		 VALUES ($1, $2, $3, $4, $5, $6)
		 ON CONFLICT (id)
		 DO UPDATE SET sender_id = EXCLUDED.sender_id, recipient_id = EXCLUDED.recipient_id, 
		               content = EXCLUDED.content, status = EXCLUDED.status, ts = EXCLUDED.ts`,
		msg.Id, msg.SenderId, msg.RecipientId, msg.Content, msg.Status, msg.Ts,
	)
	return saveErr
}

func (rep *RepositoryService) GetMessageById(messageId string) (*Message, error) {
	var msg Message
	dbErr := rep.db.QueryRow(`
		SELECT id, sender_id, recipient_id, content, status, ts 
		FROM messages 
		WHERE id = $1`, messageId).Scan(
		&msg.Id, &msg.SenderId, &msg.RecipientId, &msg.Content, &msg.Status, &msg.Ts,
	)

	if dbErr != nil {
		return nil, dbErr
	}

	return &msg, nil
}

func (rep *RepositoryService) UpdateMessage(messageId string, status MessageStatus) error {
	_, dbErr := rep.db.Exec(
		`UPDATE messages 
		 SET status = $1
		 WHERE id = $2`,
		status, messageId,
	)
	if dbErr != nil {
		return dbErr
	}
	return nil
}

func (rep *RepositoryService) GetMessages(status MessageStatus) ([]Message, error) {
	var rows *sql.Rows
	var dbErr error

	if status == StatusAll {
		rows, dbErr = rep.db.Query(`
			SELECT 
				id,
				sender_id,
				recipient_id,
				content,
				status, 
				ts 
			FROM messages 
			ORDER BY ts DESC`)
	} else {
		rows, dbErr = rep.db.Query(`
			SELECT 
				id,
				sender_id,
				recipient_id,
				content,
				status, 
				ts 
			FROM messages 
			WHERE status = $1 
			ORDER BY ts DESC`, status)
	}

	if dbErr != nil {
		return nil, dbErr
	}

	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		cursorErr := rows.Scan(&msg.Id, &msg.SenderId, &msg.RecipientId, &msg.Content, &msg.Status, &msg.Ts)
		if cursorErr != nil {
			return nil, cursorErr
		}
		messages = append(messages, msg)
	}

	if execErr := rows.Err(); execErr != nil {
		return nil, execErr
	}

	return messages, nil
}
