package app

import (
	"chat-app/internal/config"
	"chat-app/internal/database"
	"chat-app/internal/handler"
	"chat-app/internal/repository"
	"chat-app/internal/scheduler"
	"chat-app/internal/service"
	"chat-app/internal/websocket"
	"go.uber.org/fx"
)

var Modules = fx.Options(
	fx.Provide(
		config.LoadConfig,
		database.Connect,
		repository.NewChatRepository,
		service.NewChatService,
		websocket.NewHub,
		websocket.NewUpgrader,
		handler.NewChatHandler,
		scheduler.NewScheduler,
	),
	fx.Invoke(
		handler.RegisterRoutes,
		func(s *scheduler.Scheduler) {
			s.Start()
		},
	),
)
