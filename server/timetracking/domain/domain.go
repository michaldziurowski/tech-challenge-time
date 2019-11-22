package domain

import "time"

type Session struct {
	UserId    string
	SessionId int64
	Name      string
	IsOpen    bool
}

type SessionAggregate struct {
	Session
	Duration time.Duration
}

type EventType int

const (
	STARTSESSION EventType = iota
	STOPSESSION
)

type SessionEvent struct {
	Type      EventType
	UserId    string
	SessionId int64
	Time      time.Time
}
