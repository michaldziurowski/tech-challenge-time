package usecases

import (
	"time"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/domain"
)

type EventStore interface {
	AddEvent(event domain.SessionEvent) error
	GetEventsByRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error)
}

type Repository interface {
	AddSession(session domain.Session) (int64, error)
	GetSession(sessionId int64) (domain.Session, error)
	SetSessionName(sessionId int64, name string) error
	ToggleSessionState(sessionId int64) error
}

type DateProvider interface {
	GetCurrent() time.Time
}
