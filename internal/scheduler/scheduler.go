package scheduler

import (
	"chat-app/internal/repository"
	"chat-app/internal/websocket"
	"log"
	"time"
)

type Scheduler struct {
	repo repository.ChatRepository
	hub  *websocket.Hub
}

func NewScheduler(repo repository.ChatRepository, hub *websocket.Hub) *Scheduler {
	return &Scheduler{
		repo: repo,
		hub:  hub,
	}
}

func (s *Scheduler) Start() {
	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			s.run()
		}
	}()
}

func (s *Scheduler) run() {
	messages, err := s.repo.GetScheduledMessages()
	if err != nil {
		log.Println("Scheduler error fetching messages:", err)
		return
	}

	for _, msg := range messages {
		if s.hub != nil {
			s.hub.Broadcast([]byte(msg.Content))
		}
		log.Printf("Scheduled message delivered: %s\n", msg.Content)
	}
}
