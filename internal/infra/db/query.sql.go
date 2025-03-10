// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: query.sql

package db

import (
	"context"
	"time"
)

const addMessage = `-- name: AddMessage :exec
INSERT INTO messages (id, chat_id, role, content, tokens, model, erased, order_msg, created_at) VALUES(?,?,?,?,?,?,?,?,?)
`

type AddMessageParams struct {
	ID        string
	ChatID    string
	Role      string
	Content   string
	Tokens    int16
	Model     string
	Erased    bool
	OrderMsg  int16
	CreatedAt time.Time
}

func (q *Queries) AddMessage(ctx context.Context, arg AddMessageParams) error {
	_, err := q.db.ExecContext(ctx, addMessage,
		arg.ID,
		arg.ChatID,
		arg.Role,
		arg.Content,
		arg.Tokens,
		arg.Model,
		arg.Erased,
		arg.OrderMsg,
		arg.CreatedAt,
	)
	return err
}

const createChat = `-- name: CreateChat :exec
INSERT INTO chats 
    (id, user_id, initial_message_id, status, token_usage, model, model_max_tokens,temperature, top_p, n, stop, max_tokens, presence_penalty, frequency_penalty, created_at, updated_at)
    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
`

type CreateChatParams struct {
	ID               string
	UserID           string
	InitialMessageID string
	Status           string
	TokenUsage       int16
	Model            string
	ModelMaxTokens   int16
	Temperature      float32
	TopP             float32
	N                int16
	Stop             string
	MaxTokens        int16
	PresencePenalty  float32
	FrequencyPenalty float32
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) error {
	_, err := q.db.ExecContext(ctx, createChat,
		arg.ID,
		arg.UserID,
		arg.InitialMessageID,
		arg.Status,
		arg.TokenUsage,
		arg.Model,
		arg.ModelMaxTokens,
		arg.Temperature,
		arg.TopP,
		arg.N,
		arg.Stop,
		arg.MaxTokens,
		arg.PresencePenalty,
		arg.FrequencyPenalty,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	return err
}

const deleteChatMessages = `-- name: DeleteChatMessages :exec
DELETE FROM messages WHERE chat_id = ?
`

func (q *Queries) DeleteChatMessages(ctx context.Context, chatID string) error {
	_, err := q.db.ExecContext(ctx, deleteChatMessages, chatID)
	return err
}

const deleteErasedChatMessages = `-- name: DeleteErasedChatMessages :exec
DELETE FROM messages WHERE erased=1 and chat_id = ?
`

func (q *Queries) DeleteErasedChatMessages(ctx context.Context, chatID string) error {
	_, err := q.db.ExecContext(ctx, deleteErasedChatMessages, chatID)
	return err
}

const findChatByID = `-- name: FindChatByID :one
SELECT id, user_id, initial_message_id, status, token_usage, model, model_max_tokens, temperature, top_p, n, stop, max_tokens, presence_penalty, frequency_penalty, created_at, updated_at FROM chats WHERE id = ?
`

func (q *Queries) FindChatByID(ctx context.Context, id string) (Chat, error) {
	row := q.db.QueryRowContext(ctx, findChatByID, id)
	var i Chat
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.InitialMessageID,
		&i.Status,
		&i.TokenUsage,
		&i.Model,
		&i.ModelMaxTokens,
		&i.Temperature,
		&i.TopP,
		&i.N,
		&i.Stop,
		&i.MaxTokens,
		&i.PresencePenalty,
		&i.FrequencyPenalty,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const findErasedMessagesByChatID = `-- name: FindErasedMessagesByChatID :many
SELECT id, chat_id, role, content, tokens, model, erased, order_msg, created_at FROM messages WHERE erased=1 and chat_id = ? order by order_msg asc
`

func (q *Queries) FindErasedMessagesByChatID(ctx context.Context, chatID string) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, findErasedMessagesByChatID, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.ChatID,
			&i.Role,
			&i.Content,
			&i.Tokens,
			&i.Model,
			&i.Erased,
			&i.OrderMsg,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const findMessagesByChatID = `-- name: FindMessagesByChatID :many
SELECT id, chat_id, role, content, tokens, model, erased, order_msg, created_at FROM messages WHERE erased=0 and chat_id = ? order by order_msg asc
`

func (q *Queries) FindMessagesByChatID(ctx context.Context, chatID string) ([]Message, error) {
	rows, err := q.db.QueryContext(ctx, findMessagesByChatID, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.ChatID,
			&i.Role,
			&i.Content,
			&i.Tokens,
			&i.Model,
			&i.Erased,
			&i.OrderMsg,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const saveChat = `-- name: SaveChat :exec
UPDATE chats SET user_id = ?, initial_message_id = ?, status = ?, token_usage = ?, model = ?, model_max_tokens=?, temperature = ?, top_p = ?, n = ?, stop = ?, max_tokens = ?, presence_penalty = ?, frequency_penalty = ?, updated_at = ? WHERE id = ?
`

type SaveChatParams struct {
	UserID           string
	InitialMessageID string
	Status           string
	TokenUsage       int16
	Model            string
	ModelMaxTokens   int16
	Temperature      float32
	TopP             float32
	N                int16
	Stop             string
	MaxTokens        int16
	PresencePenalty  float32
	FrequencyPenalty float32
	UpdatedAt        time.Time
	ID               string
}

func (q *Queries) SaveChat(ctx context.Context, arg SaveChatParams) error {
	_, err := q.db.ExecContext(ctx, saveChat,
		arg.UserID,
		arg.InitialMessageID,
		arg.Status,
		arg.TokenUsage,
		arg.Model,
		arg.ModelMaxTokens,
		arg.Temperature,
		arg.TopP,
		arg.N,
		arg.Stop,
		arg.MaxTokens,
		arg.PresencePenalty,
		arg.FrequencyPenalty,
		arg.UpdatedAt,
		arg.ID,
	)
	return err
}
