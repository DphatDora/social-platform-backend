package service

import (
	"log"
	"social-platform-backend/internal/interface/dto/response"
	"sync"
)

// SSEClient represents a connected SSE client
type SSEClient struct {
	UserID  uint64
	Channel chan *response.SSEEvent
}

type SSEService struct {
	clients    map[uint64][]*SSEClient // list of clients (userID)
	clientsMux sync.RWMutex
}

func NewSSEService() *SSEService {
	return &SSEService{
		clients: make(map[uint64][]*SSEClient),
	}
}

func (s *SSEService) RegisterClient(userID uint64) *SSEClient {
	client := &SSEClient{
		UserID:  userID,
		Channel: make(chan *response.SSEEvent, 10),
	}

	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	s.clients[userID] = append(s.clients[userID], client)
	log.Printf("[Info] SSE client registered for user %d. Total clients: %d", userID, len(s.clients[userID]))

	return client
}

func (s *SSEService) UnregisterClient(userID uint64, client *SSEClient) {
	s.clientsMux.Lock()
	defer s.clientsMux.Unlock()

	clients := s.clients[userID]
	for i, c := range clients {
		if c == client {
			// Remove client from slice
			s.clients[userID] = append(clients[:i], clients[i+1:]...)
			close(client.Channel)
			log.Printf("[Info] SSE client unregistered for user %d. Remaining clients: %d", userID, len(s.clients[userID]))
			break
		}
	}

	// Clean up empty user entries
	if len(s.clients[userID]) == 0 {
		delete(s.clients, userID)
	}
}

// BroadcastToUser sends event to all SSE clients of a user
func (s *SSEService) BroadcastToUser(userID uint64, event *response.SSEEvent) {
	s.clientsMux.RLock()
	defer s.clientsMux.RUnlock()

	clients := s.clients[userID]
	for _, client := range clients {
		select {
		case client.Channel <- event:
			log.Printf("[Info] Event '%s' sent to user %d", event.Event, userID)
		default:
			log.Printf("[Warn] Failed to send event '%s' to user %d - channel full", event.Event, userID)
		}
	}
}
