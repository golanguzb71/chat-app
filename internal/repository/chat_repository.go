package repository

import (
	"chat-app/internal/model"
	"gorm.io/gorm"
	"time"
)

type ChatRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	BlockUser(blocked *model.BlockedUser) error
	IsBlocked(senderID, receiverID uint) (bool, error)

	CreateGroup(group *model.Group) error
	AddUserToGroup(groupID, userID uint) error
	GetGroupMembers(groupID uint) ([]model.User, error)

	CreateMessage(msg *model.Message) error
	GetMessagesByGroup(groupID uint) ([]model.Message, error)
	GetScheduledMessages() ([]model.Message, error)
}

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) ChatRepository {
	return &chatRepository{db: db}
}
