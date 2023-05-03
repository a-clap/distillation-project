/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"

	"github.com/a-clap/embedded/pkg/ds18b20"
	"github.com/a-clap/embedded/pkg/embedded"
	"github.com/a-clap/logging"
)

var (
	ErrNoSuchID      = errors.New("doesn't exist")
	ErrNoTemps       = errors.New("temperature buffer is empty")
	ErrUnexpectedID  = errors.New("unexpected ID")
	ErrNoDSInterface = errors.New("no ds interface")
)

type DSError struct {
	ID  string `json:"ID"`
	Op  string `json:"op"`
	Err string `json:"error"`
}

func (d *DSError) Error() string {
	if d.Err == "" {
		return "<nil>"
	}
	s := d.Op
	if d.ID != "" {
		s += ":" + d.ID
	}
	s += ": " + d.Err
	return s
}

// DS access to on-board DS18B20 sensors
type DS interface {
	Get() ([]embedded.DSSensorConfig, error)
	Configure(s embedded.DSSensorConfig) (embedded.DSSensorConfig, error)
	Temperatures() ([]embedded.DSTemperature, error)
}

type dsConfigureCallback func()

// DSConfig simple wrapper for sensor configuration
type DSConfig struct {
	embedded.DSSensorConfig
	temps embedded.DSTemperature
}

// DSHandler main struct used to handle number of DS sensors
type DSHandler struct {
	DS      DS
	sensors map[string]*DSConfig
	clients []dsConfigureCallback
}

// DSTemperature - json returned from rest API
type DSTemperature struct {
	ID          string  `json:"ID"`
	Temperature float64 `json:"temperature"`
}

// NewDSHandler creates new DSHandler with provided DS interface
func NewDSHandler(ds DS) (*DSHandler, error) {
	d := &DSHandler{
		DS:      ds,
		sensors: make(map[string]*DSConfig),
		clients: make([]dsConfigureCallback, 0),
	}
	if err := d.init(); err != nil {
		return nil, err
	}
	return d, nil
}

// Update updates temperatures in sensors
func (d *DSHandler) Update() (errs []error) {
	temps, err := d.DS.Temperatures()
	if err != nil {
		errs = append(errs, &DSError{Op: "Update.Temperatures", Err: err.Error()})
		return
	}
	for _, temp := range temps {
		if len(temp.Readings) == 0 {
			continue
		}
		for _, single := range temp.Readings {
			id := single.ID
			s, ok := d.sensors[id]
			if !ok {
				errs = append(errs, &DSError{ID: id, Op: "Update.Temperatures", Err: ErrUnexpectedID.Error()})
				continue
			}
			s.temps.Readings = append(s.temps.Readings, single)
		}
	}
	return
}

// History returns historical temperatures, but it also CLEARS all history data but last
func (d *DSHandler) History() []embedded.DSTemperature {
	t := make([]embedded.DSTemperature, 0, len(d.sensors))
	for _, v := range d.sensors {
		length := len(v.temps.Readings)
		if length < 2 {
			continue
		}
		var data []ds18b20.Readings
		data, v.temps.Readings = v.temps.Readings[0:length-1], v.temps.Readings[length-1:]

		t = append(t, embedded.DSTemperature{
			Readings: data,
		})
	}
	return t
}

// Temperatures returns last read temperature for all sensors
func (d *DSHandler) Temperatures() []DSTemperature {
	t := make([]DSTemperature, 0, len(d.sensors))
	for id := range d.sensors {
		temp, err := d.Temperature(id)
		if err != nil {
			logger.Debug("error on DS temperature ", logging.String("id", id), logging.String("error", err.Error()))
			continue
		}

		t = append(t, DSTemperature{
			ID:          id,
			Temperature: temp,
		})
	}
	return t
}

// Temperature returns last read temperature
func (d *DSHandler) Temperature(id string) (float64, error) {
	ds, ok := d.sensors[id]
	if !ok {
		return 0.0, &DSError{ID: id, Op: "Temperature", Err: ErrNoSuchID.Error()}
	}

	size := len(ds.temps.Readings)
	if size == 0 {
		return 0.0, &DSError{ID: id, Op: "Temperature", Err: ErrNoTemps.Error()}
	}
	// Return last temperature
	return ds.temps.Readings[size-1].Average, nil
}

func (d *DSHandler) ConfigureSensor(cfg DSConfig) (DSConfig, error) {
	newConfig := DSConfig{}
	ds, ok := d.sensors[cfg.ID]
	if !ok {
		return newConfig, &DSError{ID: cfg.ID, Op: "ConfigureSensor", Err: ErrNoSuchID.Error()}
	}
	newCfg, err := d.DS.Configure(cfg.DSSensorConfig)
	if err != nil {
		return newConfig, &DSError{ID: cfg.ID, Op: "ConfigureSensor.Set", Err: err.Error()}
	}
	ds.DSSensorConfig = newCfg
	d.notify()
	return *ds, nil
}

func (d *DSHandler) GetConfig(id string) (DSConfig, error) {
	ds, ok := d.sensors[id]
	if !ok {
		return DSConfig{}, &DSError{ID: id, Op: "GetConfig", Err: ErrNoSuchID.Error()}
	}
	return *ds, nil
}

func (d *DSHandler) GetSensors() []DSConfig {
	s := make([]DSConfig, 0, len(d.sensors))
	for _, elem := range d.sensors {
		s = append(s, *elem)
	}
	return s
}

func (d *DSHandler) init() error {
	if d.DS == nil {
		return &DSError{Op: "init", Err: ErrNoDSInterface.Error()}
	}
	sensors, err := d.DS.Get()
	if err != nil {
		return &DSError{Op: "init.Get", Err: err.Error()}
	}

	for _, ds := range sensors {
		id := ds.ID
		cfg := &DSConfig{
			DSSensorConfig: ds,
			temps:          embedded.DSTemperature{},
		}
		d.sensors[id] = cfg
		// TODO: Should we configure them on startup?
	}
	return nil
}

func (d *DSHandler) subscribe(cb dsConfigureCallback) {
	d.clients = append(d.clients, cb)
}

func (d *DSHandler) notify() {
	for _, cb := range d.clients {
		cb()
	}
}
