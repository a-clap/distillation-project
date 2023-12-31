/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package embedded_test

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"embedded/pkg/ds18b20"
	"embedded/pkg/embedded"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type DS18B20ClientSuite struct {
	suite.Suite
}

func TestDS18B20Client(t *testing.T) {
	suite.Run(t, new(DS18B20ClientSuite))
}

func (p *DS18B20ClientSuite) SetupTest() {
	gin.DefaultWriter = io.Discard
}

func (p *DS18B20ClientSuite) Test_Temperatures() {
	t := p.Require()

	args := []struct {
		cfg      embedded.DSSensorConfig
		readings []ds18b20.Readings
	}{
		{
			cfg: embedded.DSSensorConfig{
				Enabled: false,
				SensorConfig: ds18b20.SensorConfig{
					ID:           "heyo",
					Correction:   1.0,
					Resolution:   11,
					PollInterval: 100,
					Samples:      13,
				},
			},
			readings: []ds18b20.Readings{
				{
					ID:          "heyo",
					Temperature: 13,
					Average:     1,
					Stamp:       time.Time{},
					Error:       "",
				},
			},
		}, {
			cfg: embedded.DSSensorConfig{
				Enabled: false,
				SensorConfig: ds18b20.SensorConfig{
					ID:           "heyo 13",
					Correction:   1.0,
					Resolution:   11,
					PollInterval: 100,
					Samples:      13,
				},
			},
			readings: []ds18b20.Readings{
				{
					ID:          "heyo 2",
					Temperature: 14,
					Average:     51,
					Stamp:       time.Time{},
					Error:       "",
				},
			},
		},
	}

	var sensors []embedded.DSSensor
	var readings []embedded.DSTemperature
	for _, elem := range args {
		m := new(DS18B20SensorMock)
		m.On("ID").Return(elem.cfg.ID)
		m.On("GetConfig").Return(elem.cfg.SensorConfig)
		m.On("GetReadings").Return(elem.readings)
		readings = append(readings, embedded.DSTemperature{Readings: elem.readings})
		sensors = append(sensors, m)
	}

	h, _ := embedded.NewRest("", embedded.WithDS18B20(sensors))
	srv := httptest.NewServer(h.Router)
	defer srv.Close()

	ds := embedded.NewDS18B20Client(srv.URL, 1*time.Second)
	s, err := ds.Temperatures()
	t.Nil(err)
	t.NotNil(s)
	t.ElementsMatch(readings, s)
}

func (p *DS18B20ClientSuite) Test_Configure() {
	t := p.Require()

	args := []struct {
		cfg embedded.DSSensorConfig
	}{
		{
			cfg: embedded.DSSensorConfig{
				Enabled: false,
				SensorConfig: ds18b20.SensorConfig{
					ID:           "heyo",
					Correction:   1.0,
					Resolution:   11,
					PollInterval: 100,
					Samples:      13,
				},
			},
		}, {
			cfg: embedded.DSSensorConfig{
				Enabled: false,
				SensorConfig: ds18b20.SensorConfig{
					ID:           "heyo 2",
					Correction:   1.0,
					Resolution:   11,
					PollInterval: 100,
					Samples:      13,
				},
			},
		},
	}

	var sensors []embedded.DSSensor
	var cfgs []embedded.DSSensorConfig
	var mocks []*DS18B20SensorMock
	for _, elem := range args {
		m := new(DS18B20SensorMock)
		m.On("ID").Return(elem.cfg.ID)
		m.On("GetConfig").Return(elem.cfg.SensorConfig).Once()
		cfgs = append(cfgs, elem.cfg)
		mocks = append(mocks, m)
		sensors = append(sensors, m)
	}

	h, _ := embedded.NewRest("", embedded.WithDS18B20(sensors))
	srv := httptest.NewServer(h.Router)
	defer srv.Close()

	pt := embedded.NewDS18B20Client(srv.URL, 1*time.Second)
	s, err := pt.Get()
	t.Nil(err)
	t.NotNil(s)
	t.ElementsMatch(cfgs, s)

	// Expected error - sensor doesn't exist
	_, err = pt.Configure(embedded.DSSensorConfig{})
	t.NotNil(err)
	t.ErrorContains(err, embedded.ErrNoSuchID.Error())
	t.ErrorContains(err, embedded.RoutesConfigOnewireSensor)

	// Error on set now
	errSet := errors.New("hello world")
	cfgs[0].Samples = 15
	mocks[0].On("Configure", mock.Anything).Return(errSet).Once()
	_, err = pt.Configure(cfgs[0])
	t.NotNil(err)
	t.ErrorContains(err, errSet.Error())
	//
	// All good now
	cfgs[0].Samples = 23
	mocks[0].On("Configure", mock.Anything).Return(nil)
	mocks[0].On("GetConfig").Return(cfgs[0].SensorConfig).Once()
	cfg, err := pt.Configure(cfgs[0])
	t.Nil(err)
	t.Equal(cfgs[0], cfg)
}

func (p *DS18B20ClientSuite) Test_NotImplemented() {
	t := p.Require()
	h, _ := embedded.NewRest("")
	srv := httptest.NewServer(h.Router)
	defer srv.Close()

	pt := embedded.NewDS18B20Client(srv.URL, 1*time.Second)

	s, err := pt.Get()
	t.Nil(s)
	t.NotNil(err)
	t.ErrorContains(err, embedded.ErrNotImplemented.Error())
	t.ErrorContains(err, embedded.RoutesGetOnewireSensors)

	_, err = pt.Configure(embedded.DSSensorConfig{})
	t.NotNil(err)
	t.ErrorContains(err, embedded.ErrNotImplemented.Error())
	t.ErrorContains(err, embedded.RoutesConfigOnewireSensor)

	temps, err := pt.Temperatures()
	t.Nil(temps)
	t.NotNil(err)
	t.ErrorContains(err, embedded.ErrNotImplemented.Error())
	t.ErrorContains(err, embedded.RoutesGetOnewireTemperatures)
}
