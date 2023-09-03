package backendmock

import (
	"osservice"
	"time"
)

var _ osservice.Time = (*TimeMock)(nil)

type TimeMock struct{}

func (TimeMock) Now() (time.Time, error) {
	return time.Now(), nil
}

func (TimeMock) SetNow(time.Time) error {
	return nil
}

func (TimeMock) NTP() (bool, error) {
	return false, nil
}

func (TimeMock) SetNTP(enable bool) error {
	return nil
}
