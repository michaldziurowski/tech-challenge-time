package usecases

import (
	"testing"
	"time"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/domain"
)

const (
	USER        string = "user"
	SESSIONNAME string = "session"
	SESSIONID1  int64  = 1
	SESSIONID2  int64  = 2
)

type eventStoreStub struct {
	getEventsByRangeStub func(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error)
}

func NewEventStoreStub(getEventsByRangeStub func(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error)) EventStore {
	return eventStoreStub{getEventsByRangeStub}
}

func (e eventStoreStub) AddEvent(event domain.SessionEvent) error {
	return nil
}

func (e eventStoreStub) GetEventsByRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return e.getEventsByRangeStub(userId, from, to)
}

type repositoryStub struct{}

func (r repositoryStub) AddSession(session domain.Session) (int64, error) { return 1, nil }

func (r repositoryStub) GetSession(sessionId int64) (domain.Session, error) {
	return domain.Session{
		UserId:    USER,
		SessionId: SESSIONID1,
		Name:      SESSIONNAME,
	}, nil
}

func (r repositoryStub) SetSessionName(sessionId int64, name string) error { return nil }

func (r repositoryStub) ToggleSessionState(sessionId int64) error { return nil }

type dateProviderStub struct {
	currentTime time.Time
}

func NewDateProviderStub(currentTime time.Time) DateProvider {
	return dateProviderStub{currentTime}
}

func (d dateProviderStub) GetCurrent() time.Time {
	return d.currentTime
}

func getTime(timeString string) time.Time {
	time, _ := time.Parse(time.RFC3339, timeString)
	return time
}

func getDuration(duratonString string) time.Duration {
	duration, _ := time.ParseDuration(duratonString)
	return duration
}

func multipleFinishedSessionsInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T16:04:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T16:15:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T17:00:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T18:15:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID2,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T19:00:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID2,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T20:15:05Z"),
		},
	}, nil
}

func singleFinishedSessionInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T16:04:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T16:15:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T17:00:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T18:15:05Z"),
		},
	}, nil
}

func singleClosedAtFrontSessionInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T16:15:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T17:00:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T18:15:05Z"),
		},
	}, nil
}

func singleOpenAtEndSessionInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T17:00:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T18:15:05Z"),
		},
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T19:15:05Z"),
		},
	}, nil
}

func singleWithOnlyOneClosedAtFrontSessionInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STOPSESSION,
			Time:      getTime("2019-11-20T16:15:05Z"),
		},
	}, nil
}

func singleWithOnlyOneOpenAtEndSessionInRange(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error) {
	return []domain.SessionEvent{
		domain.SessionEvent{
			UserId:    USER,
			SessionId: SESSIONID1,
			Type:      domain.STARTSESSION,
			Time:      getTime("2019-11-20T19:15:05Z"),
		},
	}, nil
}
func TestGetSessionByRangeWhenSingleSessionInRange(t *testing.T) {
	testCases := []struct {
		caseMethod              func(userId string, from time.Time, to time.Time) ([]domain.SessionEvent, error)
		expectedAggregatesCount int
		expectedDuration        time.Duration
	}{
		{singleFinishedSessionInRange, 1, getDuration("1h26m")},
		{singleClosedAtFrontSessionInRange, 1, getDuration("1h15m")},
		{singleOpenAtEndSessionInRange, 1, getDuration("5h15m")},
		{singleWithOnlyOneClosedAtFrontSessionInRange, 0, getDuration("0m")},
		{singleWithOnlyOneOpenAtEndSessionInRange, 1, getDuration("4h")},
	}

	for _, testCase := range testCases {
		esStub := NewEventStoreStub(testCase.caseMethod)
		service := NewService(esStub, repositoryStub{}, NewDateProviderStub(getTime("2019-11-20T23:15:05Z")))

		notRelevantTime := time.Now()
		aggregates, _ := service.GetSessionsByRange(USER, notRelevantTime, notRelevantTime)

		if len(aggregates) != testCase.expectedAggregatesCount {
			t.Errorf("Expected %v aggregates but got %v", testCase.expectedAggregatesCount, len(aggregates))
		}

		if len(aggregates) != 0 && aggregates[0].Duration != testCase.expectedDuration {
			t.Errorf("For session %s expected duration is %v but was %v", aggregates[0].Name, testCase.expectedDuration, aggregates[0].Duration)
		}
	}
}

func TestGetSessionByRangeWhenMultipleSessionsInRange(t *testing.T) {
	esStub := NewEventStoreStub(multipleFinishedSessionsInRange)
	service := NewService(esStub, repositoryStub{}, NewDateProviderStub(getTime("2019-11-20T23:15:05Z")))

	notRelevantTime := time.Now()
	aggregates, _ := service.GetSessionsByRange(USER, notRelevantTime, notRelevantTime)

	expectedAggregatesCount := 2
	if len(aggregates) != expectedAggregatesCount {
		t.Errorf("Expected %v aggregates but got %v", expectedAggregatesCount, len(aggregates))
	}

	firstAggregateExpectedDuration := getDuration("1h26m")
	if aggregates[0].Duration != firstAggregateExpectedDuration {
		t.Errorf("For session %s expected duration is %v but was %v", aggregates[0].Name, firstAggregateExpectedDuration, aggregates[0].Duration)
	}

	secondAggregateExpectedDuration := getDuration("1h15m")
	if aggregates[1].Duration != secondAggregateExpectedDuration {
		t.Errorf("For session %s expected duration is %v but was %v", aggregates[1].Name, secondAggregateExpectedDuration, aggregates[1].Duration)
	}
}
