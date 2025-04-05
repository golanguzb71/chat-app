package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"time"
)

type ChatService interface {
	RegisterUser(username, password string) (*model.User, error)
	BlockUser(blockerID, blockedID uint) error

	CreateGroup(name *string, memberIDs []uint) (*model.Group, error)
	SendMessage(senderID uint, groupID *uint, content string, scheduledAt *time.Time) error
	GetGroupMessages(groupID uint) ([]model.Message, error)
}

type chatService struct {
	repo repository.ChatRepository
}

func NewChatService(repo repository.ChatRepository) ChatService {
	return &chatService{repo: repo}
}
