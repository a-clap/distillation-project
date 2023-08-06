package process

import (
	"errors"
	"testing"

	"distillation/pkg/process/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSensor_TemperatureFailure(t *testing.T) {
	r := require.New(t)

	ctrl := gomock.NewController(t)
	s := mocks.NewMockSensor(ctrl)
	s.EXPECT().ID().Return("s1").AnyTimes()

	ps := newSensor(s)

	r.EqualValues(ps.ID(), "s1")
	// Return correct temperature
	correct := 12.0
	s.EXPECT().Temperature().Return(correct, nil)
	tmp, err := ps.Temperature()
	r.InDelta(correct, tmp, 0.1)
	r.Nil(err)

	// // Underlying mock returns error
	retErr := errors.New("")
	for i := 0; i < maxFailures; i++ {
		s.EXPECT().Temperature().Return(correct, retErr)
		tmp, err = ps.Temperature()
		// Should return the same temperature
		r.InDelta(correct, tmp, 0.1)
		r.Nil(err)
	}
	// Now we should receive error
	s.EXPECT().Temperature().Return(correct, retErr)
	_, err = ps.Temperature()
	// Without error
	r.NotNil(err)
	r.ErrorIs(err, ErrTooManyErrorOnTemperature)
}
