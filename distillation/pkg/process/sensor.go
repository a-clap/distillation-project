package process

import (
	"fmt"
)

const maxFailures = 3

type Sensor interface {
	ID() string                    // ID returns unique ID of interface
	Temperature() (float64, error) // Temperature returns latest temperature read from sensor
}

type sensor struct {
	Sensor   Sensor
	temp     float64
	failures int
}

func newSensor(s Sensor) *sensor {
	return &sensor{
		Sensor:   s,
		temp:     0,
		failures: 0,
	}
}

func (s *sensor) ID() string {
	return s.Sensor.ID()
}

func (s *sensor) Temperature() (float64, error) {
	tmp, err := s.Sensor.Temperature()
	if err != nil {
		s.failures++
		// Something must be wrong
		if s.failures > maxFailures {
			return 0, fmt.Errorf("%w: on ID: %v", ErrTooManyErrorOnTemperature, s.ID())
		}
	} else {
		// All good, update temperature and reset failures
		s.failures = 0
		s.temp = tmp
	}
	return s.temp, nil
}
