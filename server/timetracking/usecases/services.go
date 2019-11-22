package usecases

import (
	"fmt"
	"time"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/domain"
)

type Service interface {
	StartSession(userId string, name string, startedAt time.Time) (int64, error)
	StopSession(userId string, sessionId int64, stoppedAt time.Time) error
	ResumeSession(userId string, sessionId int64, resumedAt time.Time) error
	SetSessionName(userId string, sessionId int64, name string) error
	GetSessionsByRange(userId string, from time.Time, to time.Time) ([]domain.SessionAggregate, error)
}

type service struct {
	eventStore   EventStore
	repository   Repository
	dateProvider DateProvider
}

func NewService(e EventStore, r Repository, d DateProvider) Service {
	return service{e, r, d}
}

func (s service) getSessionForUser(userId string, sessionId int64) (domain.Session, error) {
	session, err := s.repository.GetSession(sessionId)

	if err != nil {
		return domain.Session{}, err
	}

	if session.UserId != userId {
		return domain.Session{}, fmt.Errorf("Session doesnt belong to user.")
	}

	return session, nil
}

func (s service) StartSession(userId string, name string, startedAt time.Time) (int64, error) {
	session := domain.Session{
		UserId: userId,
		Name:   name,
		IsOpen: true,
	}

	sessionId, err := s.repository.AddSession(session)

	if err != nil {
		return -1, err
	}

	startEvent := domain.SessionEvent{
		Type:      domain.STARTSESSION,
		UserId:    userId,
		SessionId: sessionId,
		Time:      startedAt,
	}

	err = s.eventStore.AddEvent(startEvent)

	if err != nil {
		return -1, err
	}

	return sessionId, nil
}

func (s service) StopSession(userId string, sessionId int64, stoppedAt time.Time) error {
	_, err := s.getSessionForUser(userId, sessionId)

	if err != nil {
		return err
	}

	err = s.repository.ToggleSessionState(sessionId)

	if err != nil {
		return err
	}

	stopEvent := domain.SessionEvent{
		Type:      domain.STOPSESSION,
		UserId:    userId,
		SessionId: sessionId,
		Time:      stoppedAt,
	}

	err = s.eventStore.AddEvent(stopEvent)

	if err != nil {
		return err
	}

	return nil
}

func (s service) ResumeSession(userId string, sessionId int64, resumedAt time.Time) error {
	_, err := s.getSessionForUser(userId, sessionId)

	if err != nil {
		return err
	}

	err = s.repository.ToggleSessionState(sessionId)

	if err != nil {
		return err
	}

	startEvent := domain.SessionEvent{
		Type:      domain.STARTSESSION,
		UserId:    userId,
		SessionId: sessionId,
		Time:      resumedAt,
	}

	err = s.eventStore.AddEvent(startEvent)

	if err != nil {
		return err
	}

	return nil
}

func (s service) SetSessionName(userId string, sessionId int64, name string) error {
	_, err := s.getSessionForUser(userId, sessionId)

	if err != nil {
		return err
	}

	err = s.repository.SetSessionName(sessionId, name)

	if err != nil {
		return err
	}

	return nil
}

type sessionDuration struct {
	startTime time.Time
	duration  time.Duration
}

func (s service) GetSessionsByRange(userId string, from time.Time, to time.Time) ([]domain.SessionAggregate, error) {
	events, err := s.eventStore.GetEventsByRange(userId, from, to)
	if err != nil {
		return nil, err
	}

	durations := make(map[int64]*sessionDuration)

	for idx, event := range events {
		if sDuration, exist := durations[event.SessionId]; exist {
			switch event.Type {
			case domain.STOPSESSION:
				sDuration.duration += event.Time.Sub(sDuration.startTime)
			case domain.STARTSESSION:
				if idx == len(events)-1 {
					sDuration.duration += s.dateProvider.GetCurrent().Sub(event.Time)
				} else {
					sDuration.startTime = event.Time
				}
			}
		} else if event.Type == domain.STARTSESSION {
			durations[event.SessionId] = &sessionDuration{startTime: event.Time}
			if idx == len(events)-1 {
				durations[event.SessionId].duration += s.dateProvider.GetCurrent().Sub(event.Time)
			}
		}
	}

	aggregates := make([]domain.SessionAggregate, 0, len(durations))

	for sessionId, sDuration := range durations {
		session, _ := s.repository.GetSession(sessionId)

		aggregate := domain.SessionAggregate{
			Session:  session,
			Duration: sDuration.duration,
		}

		aggregates = append(aggregates, aggregate)
	}

	return aggregates, nil
}
