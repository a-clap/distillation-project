/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation_test

import (
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/a-clap/distillation/pkg/distillation"
	"github.com/a-clap/distillation/pkg/process"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ProcessClientSuite struct {
	suite.Suite
}

func TestProcessClient(t *testing.T) {
	suite.Run(t, new(ProcessClientSuite))
}

func (p *ProcessClientSuite) SetupTest() {
	gin.DefaultWriter = io.Discard
}

type ProcessHeaterMock struct {
	mock.Mock
}
type ProcessSensorMock struct {
	mock.Mock
}

func (s *ProcessSensorMock) ID() string {
	return s.Called().String(0)
}

func (s *ProcessSensorMock) Temperature() (float64, error) {
	args := s.Called()
	return args.Get(0).(float64), args.Error(1)
}

func (p *ProcessHeaterMock) ID() string {
	return p.Called().String(0)
}

func (p *ProcessHeaterMock) SetPower(pwr int) error {
	return p.Called(pwr).Error(0)
}

func (p *ProcessClientSuite) Test_Config() {
	t := p.Require()
	h, _ := distillation.NewRest("")

	sensorMock := new(ProcessSensorMock)
	sensorMock.On("ID").Return("s1")
	sensors := []process.Sensor{sensorMock}

	heaterMock := new(ProcessHeaterMock)
	heaterMock.On("ID").Return("h1")
	heaters := []process.Heater{heaterMock}

	h.Process.UpdateHeaters(heaters)
	h.Process.UpdateSensors(sensors)
	srv := httptest.NewServer(h)
	defer srv.Close()

	ps := distillation.NewProcessClient(srv.URL, 1*time.Second)

	// Get anything
	cfg, err := ps.GetPhaseConfig(0)
	t.Nil(err)
	t.NotNil(cfg)

	// Correct config
	cfg.Heaters = []process.HeaterPhaseConfig{
		{
			ID:    "h1",
			Power: 0,
		},
	}
	cfg.Next.TimeLeft = 1

	newCfg, err := ps.ConfigurePhase(0, cfg)
	t.Nil(err)
	t.NotNil(newCfg)
	t.Equal(cfg, newCfg)

	// Ask for not existing phase
	newCfg, err = ps.ConfigurePhase(3, cfg)
	t.NotNil(err)
	t.ErrorContains(err, distillation.RoutesProcessConfigPhase)
	t.ErrorContains(err, process.ErrNoSuchPhase.Error())

}

func (p *ProcessClientSuite) Test_PhaseCount() {
	t := p.Require()
	h, _ := distillation.NewRest("")
	srv := httptest.NewServer(h)
	defer srv.Close()

	ps := distillation.NewProcessClient(srv.URL, 1*time.Second)
	s, err := ps.GetPhaseCount()
	t.Nil(err)
	t.NotNil(s)
	// Initial value
	t.EqualValues(3, s.PhaseNumber)

	// Good
	s.PhaseNumber = 5
	s, err = ps.ConfigurePhaseCount(s)
	t.Nil(err)
	t.EqualValues(5, s.PhaseNumber)

	// So there should be still 5 phases
	s, err = ps.GetPhaseCount()
	t.Nil(err)
	t.NotNil(s)
	t.EqualValues(5, s.PhaseNumber)
}
