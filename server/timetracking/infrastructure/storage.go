package infrastructure

import (
	"fmt"
	"sync"
	"time"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/domain"
)

type InMemoryStorage struct {
	mu       sync.Mutex
	sessions []domain.Session
	events   []domain.SessionEvent
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		sessions: make([]domain.Session, 0, 10),
		events:   make([]domain.SessionEvent, 0, 10),
	}
}

func (storage *InMemoryStorage) AddSession(session domain.Session) (int64, error) {
	storage.mu.Lock()
	session.SessionId = int64(len(storage.sessions) + 1)
	storage.sessions = append(storage.sessions, session)
	storage.mu.Unlock()
	return session.SessionId, nil
}

func (storage *InMemoryStorage) GetSession(sessionId int64) (domain.Session, error) {
	for i := 0; i < len(storage.sessions); i++ {
		if storage.sessions[i].SessionId == sessionId {
			return storage.sessions[i], nil
		}
	}

	return domain.Session{}, fmt.Errorf("Session not found")
}

func (storage *InMemoryStorage) SetSessionName(sessionId int64, name string) error {
	found := false
	for i := 0; i < len(storage.sessions) && !found; i++ {
		if storage.sessions[i].SessionId == sessionId {
			storage.sessions[i].Name = name
			found = true
		}
	}

	return nil
}

func (storage *InMemoryStorage) ToggleSessionState(sessionId int64) error {
	found := false
	for i := 0; i < len(storage.sessions) && !found; i++ {
		if storage.sessions[i].SessionId == sessionId {
			storage.sessions[i].IsOpen = !storage.sessions[i].IsOpen
			found = true
		}
	}
	return nil
}

func (storage *InMemoryStorage) AddEvent(event domain.SessionEvent) error {
	storage.mu.Lock()
	storage.events = append(storage.events, event)
	storage.mu.Unlock()
	return nil
}

func (storage *InMemoryStorage) GetEventsByRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	events := make([]domain.SessionEvent, 0, 10)

	for _, event := range storage.events {
		if event.UserId == userId && (event.Time.After(from) || event.Time.Equal(from)) && (event.Time.Before(to) || event.Time.Equal(to)) {
			events = append(events, event)
		}
	}

	return events, nil
}
