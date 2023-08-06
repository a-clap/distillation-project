/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"
	"time"

	"embedded/pkg/embedded"
	"embedded/pkg/max31865"
	"github.com/a-clap/logging"
)

var (
	ErrNoPTInterface = errors.New("no pt interface")
)

type PTError struct {
	ID  string `json:"ID"`
	Op  string `json:"op"`
	Err string `json:"error"`
}

func (e *PTError) Error() string {
	if e.Err == "" {
		return "<nil>"
	}
	s := e.Op
	if e.ID != "" {
		s += ":" + e.ID
	}
	s += ": " + e.Err
	return s
}

// PT access to on-board PT100 sensors
type PT interface {
	Get() ([]embedded.PTSensorConfig, error)
	Configure(s embedded.PTSensorConfig) (embedded.PTSensorConfig, error)
	Temperatures() ([]embedded.PTTemperature, error)
}

// PTConfig simple wrapper for sensor configuration
type PTConfig struct {
	embedded.PTSensorConfig
	temps embedded.PTTemperature
}

// PTHandler main struct used to handle number of PT sensors
type PTHandler struct {
	PT           PT
	sensors      map[string]*PTConfig
	pollInterval time.Duration
	clients      []ptCallback
}

// Temperature - json returned from rest API
type Temperature struct {
	ID          string    `json:"ID"`
	Temperature float64   `json:"temperature"`
	Stamp       int64     `json:"unix_seconds"`
	Error       ErrorCode `json:"error_code"`
}

type ptCallback func()

// NewPTHandler creates new PTHandler with provided PT interface
func NewPTHandler(pt PT) (*PTHandler, error) {
	d := &PTHandler{
		PT:      pt,
		sensors: make(map[string]*PTConfig),
		clients: make([]ptCallback, 0),
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Update updates temperatures in sensors
func (p *PTHandler) Update() (errs []error) {
	temps, err := p.PT.Temperatures()
	if err != nil {
		errs = append(errs, &PTError{Op: "Update.Temperatures", Err: err.Error()})
		return
	}
	for _, temp := range temps {
		if len(temp.Readings) == 0 {
			continue
		}
		for _, single := range temp.Readings {
			id := single.ID
			s, ok := p.sensors[id]
			if !ok {
				errs = append(errs, &PTError{ID: id, Op: "Update.Temperatures", Err: ErrUnexpectedID.Error()})
				continue
			}
			s.temps.Readings = append(s.temps.Readings, single)
		}
	}
	return
}

// History returns historical temperatures, but it also CLEARS all history data but last
func (p *PTHandler) History() []embedded.PTTemperature {
	t := make([]embedded.PTTemperature, 0, len(p.sensors))
	for _, v := range p.sensors {
		length := len(v.temps.Readings)
		if length < 2 {
			continue
		}
		var data []max31865.Readings
		data, v.temps.Readings = v.temps.Readings[0:length-1], v.temps.Readings[length-1:]

		t = append(t, embedded.PTTemperature{
			Readings: data,
		})
	}
	return t
}

// Temperatures returns last read temperature for all sensors
func (p *PTHandler) Temperatures() []Temperature {
	t := make([]Temperature, 0, len(p.sensors))
	for id := range p.sensors {
		t = append(t, p.Temperature(id))
	}
	return t
}

// Temperature returns last read temperature
func (p *PTHandler) Temperature(id string) Temperature {
	t := Temperature{
		ID:          id,
		Temperature: 0,
		Stamp:       0,
		Error:       0,
	}

	pt, ok := p.sensors[id]
	if !ok {

		t.Error = ErrorCodeWrongID
		return t
	}

	size := len(pt.temps.Readings)
	if size == 0 {
		logger.Error("PT.Temperature error", logging.String("id", id), logging.String("error", ErrNoTemps.Error()))
		t.Error = ErrorCodeEmptyBuffer
		return t
	}
	temp := pt.temps.Readings[size-1]
	if temp.Error != "" {
		t.Error = ErrorCodeInternal
	} else {
		t.Stamp = temp.Stamp.Unix()
		t.Temperature = temp.Average
	}

	return t
}

func (p *PTHandler) Configure(cfg PTConfig) (PTConfig, error) {
	c := PTConfig{}
	pt, ok := p.sensors[cfg.ID]
	if !ok {
		return c, &PTError{ID: cfg.ID, Op: "Configure", Err: ErrNoSuchID.Error()}
	}
	newCfg, err := p.PT.Configure(cfg.PTSensorConfig)
	if err != nil {
		return c, &PTError{ID: cfg.ID, Op: "Configure.Set", Err: err.Error()}
	}
	pt.PTSensorConfig = newCfg
	p.notify()
	return *pt, nil
}

func (p *PTHandler) GetConfig(id string) (PTConfig, error) {
	pt, ok := p.sensors[id]
	if !ok {
		return PTConfig{}, &PTError{ID: id, Op: "GetConfig", Err: ErrNoSuchID.Error()}
	}
	return *pt, nil
}

func (p *PTHandler) GetSensors() []PTConfig {
	s := make([]PTConfig, 0, len(p.sensors))
	for _, elem := range p.sensors {
		s = append(s, *elem)
	}
	return s
}

func (p *PTHandler) init() error {

	if p.PT == nil {
		return &PTError{Op: "init", Err: ErrNoPTInterface.Error()}
	}
	sensors, err := p.PT.Get()
	if err != nil {
		return &PTError{Op: "init.Get", Err: err.Error()}
	}

	for _, pt := range sensors {
		id := pt.ID
		cfg := &PTConfig{
			PTSensorConfig: pt,
			temps:          embedded.PTTemperature{},
		}
		p.sensors[id] = cfg
	}
	return nil
}

func (p *PTHandler) subscribe(cb ptCallback) {
	p.clients = append(p.clients, cb)
}
func (p *PTHandler) notify() {
	for _, client := range p.clients {
		client()
	}
}
