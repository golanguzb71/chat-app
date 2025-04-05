package service

import (
	"chat-app/internal/model"
	"chat-app/internal/repository"
	"errors"
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

func (s *chatService) RegisterUser(username, password string) (*model.User, error) {
	existing, _ := s.repo.GetUserByUsername(username)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("user already exists")
	}

	user := &model.User{
		Username: username,
		Password: password,
	}
	err := s.repo.CreateUser(user)
	return user, err
}

func (s *chatService) BlockUser(blockerID, blockedID uint) error {
	blocked := &model.BlockedUser{
		BlockerID: blockerID,
		BlockedID: blockedID,
		CreatedAt: time.Now(),
	}
	return s.repo.BlockUser(blocked)
}

func (s *chatService) CreateGroup(name *string, memberIDs []uint) (*model.Group, error) {
	group := &model.Group{
		Name:      name,
		CreatedAt: time.Now(),
	}
	if err := s.repo.CreateGroup(group); err != nil {
		return nil, err
	}

	for _, uid := range memberIDs {
		if err := s.repo.AddUserToGroup(group.ID, uid); err != nil {
			return nil, err
		}
	}

	return group, nil
}

func (s *chatService) SendMessage(senderID uint, groupID *uint, content string, scheduledAt *time.Time) error {
	if groupID != nil {
		members, err := s.repo.GetGroupMembers(*groupID)
		if err != nil {
			return err
		}
		for _, member := range members {
			blocked, err := s.repo.IsBlocked(senderID, member.ID)
			if err != nil {
				return err
			}
			if blocked {
				return errors.New("you are blocked by someone in this group")
			}
		}
	}

	msg := &model.Message{
		SenderID:    senderID,
		GroupID:     groupID,
		Content:     content,
		ScheduledAt: scheduledAt,
		CreatedAt:   time.Now(),
	}
	return s.repo.CreateMessage(msg)
}

func (s *chatService) GetGroupMessages(groupID uint) ([]model.Message, error) {
	return s.repo.GetMessagesByGroup(groupID)
}
