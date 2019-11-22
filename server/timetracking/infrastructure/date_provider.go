package infrastructure

import (
	"time"

	"github.com/michaldziurowski/tech-challenge-time/server/timetracking/usecases"
)

type dateProvider struct{}

func NewDateProvider() usecases.DateProvider {
	return dateProvider{}
}

func (dateProvider) GetCurrent() time.Time {
	return time.Now().UTC()
}
