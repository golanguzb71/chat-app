package handler

import (
	"chat-app/internal/service"
	"chat-app/internal/websocket"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type ChatHandler struct {
	service  service.ChatService
	upgrader *websocket.Upgrader
}

func NewChatHandler(service service.ChatService, upgrader *websocket.Upgrader) *ChatHandler {
	return &ChatHandler{
		service:  service,
		upgrader: upgrader,
	}
}

func RegisterRoutes(handler *ChatHandler) {
	go handler.upgrader.Hub.Run()

	r := chi.NewRouter()

	r.Post("/register", handler.RegisterUser)
	r.Post("/block", handler.BlockUser)
	r.Post("/group", handler.CreateGroup)
	r.Post("/send", handler.SendMessage)
	r.Get("/messages/{groupID}", handler.GetMessages)
	r.Get("/ws", handler.WebSocketEndpoint)

	http.ListenAndServe(":8080", r)
}

func (h *ChatHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	user, err := h.service.RegisterUser(req.Username, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *ChatHandler) BlockUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BlockerID uint `json:"blocker_id"`
		BlockedID uint `json:"blocked_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err := h.service.BlockUser(req.BlockerID, req.BlockedID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name      *string `json:"name"`
		MemberIDs []uint  `json:"member_ids"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	group, err := h.service.CreateGroup(req.Name, req.MemberIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(group)
}

func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SenderID    uint    `json:"sender_id"`
		GroupID     *uint   `json:"group_id"`
		Content     string  `json:"content"`
		ScheduledAt *string `json:"scheduled_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var parsedTime *time.Time
	if req.ScheduledAt != nil {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err == nil {
			parsedTime = &t
		}
	}

	err := h.service.SendMessage(req.SenderID, req.GroupID, req.Content, parsedTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	groupIDStr := chi.URLParam(r, "groupID")
	groupID, err := strconv.Atoi(groupIDStr)
	if err != nil {
		http.Error(w, "invalid group ID", http.StatusBadRequest)
		return
	}

	messages, err := h.service.GetGroupMessages(uint(groupID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(messages)
}

func (h *ChatHandler) WebSocketEndpoint(w http.ResponseWriter, r *http.Request) {
	h.upgrader.HandleWebSocket(w, r)
}
