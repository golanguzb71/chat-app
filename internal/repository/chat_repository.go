package repository

import (
	"chat-app/internal/model"
	"errors"
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

func (r *chatRepository) CreateUser(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *chatRepository) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	return &user, err
}

func (r *chatRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *chatRepository) BlockUser(blocked *model.BlockedUser) error {
	return r.db.Create(blocked).Error
}

func (r *chatRepository) IsBlocked(senderID, receiverID uint) (bool, error) {
	var record model.BlockedUser
	err := r.db.Where("blocker_id = ? AND blocked_id = ?", receiverID, senderID).First(&record).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return err == nil, err
}

func (r *chatRepository) CreateGroup(group *model.Group) error {
	return r.db.Create(group).Error
}

func (r *chatRepository) AddUserToGroup(groupID, userID uint) error {
	member := model.GroupMember{
		GroupID:  groupID,
		UserID:   userID,
		JoinedAt: time.Now(),
	}
	return r.db.Create(&member).Error
}

func (r *chatRepository) GetGroupMembers(groupID uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Table("users").
		Joins("JOIN group_members ON users.id = group_members.user_id").
		Where("group_members.group_id = ?", groupID).
		Scan(&users).Error
	return users, err
}

func (r *chatRepository) CreateMessage(msg *model.Message) error {
	return r.db.Create(msg).Error
}

func (r *chatRepository) GetMessagesByGroup(groupID uint) ([]model.Message, error) {
	var messages []model.Message
	err := r.db.
		Where("group_id = ?", groupID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

func (r *chatRepository) GetScheduledMessages() ([]model.Message, error) {
	var messages []model.Message
	now := time.Now()
	err := r.db.
		Where("scheduled_at IS NOT NULL AND scheduled_at <= ?", now).
		Find(&messages).Error
	return messages, err
}
