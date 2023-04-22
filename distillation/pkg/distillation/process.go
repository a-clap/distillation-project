/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"github.com/a-clap/distillation/pkg/distillation/process"
	"github.com/a-clap/embedded/pkg/embedded"
)

// ProcessPhaseCount is JSON wrapper for process.PhaseNumber
type ProcessPhaseCount struct {
	PhaseNumber int `json:"phase_number"`
}

// ProcessPhaseConfig is package wrapper for process.PhaseConfig
type ProcessPhaseConfig struct {
	process.PhaseConfig
}

// ProcessConfig serves as a way to enable or disable processing
type ProcessConfig struct {
	Enable     bool `json:"enable"`
	MoveToNext bool `json:"move_to_next"`
	Disable    bool `json:"disable"`
}

// ProcessConfigValidation returns whether Process configuration is valid
type ProcessConfigValidation struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

// ProcessStatus is just wrapper for process.Status
type ProcessStatus struct {
	process.Status
}

func (h *Distillation) configurePhase(number int, config ProcessPhaseConfig) error {
	// Update ios, if process is not running
	if h.Process.Running() == false {
		h.updateProcess()
	}
	return h.Process.ConfigurePhase(number, config.PhaseConfig)
	
}

func (h *Distillation) updateProcess() {
	h.updateHeaters()
	h.updateOutputs()
	h.updateSensors()
}

func (h *Distillation) safeUpdateSensors() {
	if h.Process.Running() == false {
		h.updateSensors()
	}
}
func (h *Distillation) safeUpdateHeaters() {
	if h.Process.Running() == false {
		h.updateHeaters()
	}
}
func (h *Distillation) safeUpdateOutputs() {
	if h.Process.Running() == false {
		h.updateOutputs()
	}
}

func (h *Distillation) updateSensors() {
	if h.DSHandler == nil && h.PTHandler == nil {
		return
	}
	
	getTempDS := func(id string) func() float64 {
		return func() float64 {
			t, err := h.DSHandler.Temperature(id)
			if err != nil {
				return -1
			}
			return t
		}
	}
	
	getTempPT := func(id string) func() float64 {
		return func() float64 {
			t, err := h.PTHandler.Temperature(id)
			if err != nil {
				return -1
			}
			return t
		}
	}
	var sensors []process.Sensor
	if h.DSHandler != nil {
		for _, ds := range h.DSHandler.GetSensors() {
			if ds.Enabled {
				s := &processSensor{
					id:      ds.ID,
					getTemp: getTempDS(ds.ID),
				}
				sensors = append(sensors, s)
			}
		}
	}
	if h.PTHandler != nil {
		for _, pt := range h.PTHandler.GetSensors() {
			if pt.Enabled {
				s := &processSensor{
					id:      pt.ID,
					getTemp: getTempPT(pt.ID),
				}
				sensors = append(sensors, s)
			}
		}
	}
	h.Process.ConfigureSensors(sensors)
	
}
func (h *Distillation) updateOutputs() {
	if h.GPIOHandler == nil {
		return
	}
	setValue := func(id string) func(value bool) error {
		return func(value bool) error {
			o, ok := h.GPIOHandler.io[id]
			if !ok {
				return ErrNoSuchID
			}
			cfg := *o
			cfg.Value = value
			_, err := h.GPIOHandler.Configure(cfg)
			return err
		}
	}
	
	var outputs []process.Output
	for _, out := range h.GPIOHandler.Config() {
		o := &processOutput{
			id:       out.ID,
			setValue: setValue(out.ID),
		}
		outputs = append(outputs, o)
	}
	h.Process.ConfigureOutputs(outputs)
	
}

func (h *Distillation) updateHeaters() {
	if h.HeatersHandler == nil {
		return
	}
	
	setPwr := func(id string) func(pwr int) error {
		return func(pwr int) error {
			cfg := HeaterConfig{
				HeaterConfig: embedded.HeaterConfig{
					ID:      id,
					Enabled: true,
					Power:   uint(pwr),
				},
			}
			_, err := h.HeatersHandler.Configure(cfg)
			return err
		}
	}
	
	var heaters []process.Heater
	for _, heater := range h.HeatersHandler.ConfigsGlobal() {
		if heater.Enabled {
			h := &processHeater{
				id:     heater.ID,
				setPwr: setPwr(heater.ID)}
			heaters = append(heaters, h)
		}
	}
	h.Process.ConfigureHeaters(heaters)
}

// processHeater fulfills process.Heater interface
type processHeater struct {
	id     string
	setPwr func(pwr int) error
}

// ID fulfills process.Heater interface
func (p *processHeater) ID() string {
	return p.id
}

// SetPower fulfills process.Heater interface
func (p *processHeater) SetPower(pwr int) error {
	return p.setPwr(pwr)
}

// processOutput fulfills process.Output interface
type processOutput struct {
	id       string
	setValue func(value bool) error
}

// ID fulfills process.Output interface
func (p *processOutput) ID() string {
	return p.id
}

// Set fulfills process.Output interface
func (p *processOutput) Set(value bool) error {
	return p.setValue(value)
}

// processSensor fulfills process.Sensor interface
type processSensor struct {
	id      string
	getTemp func() float64
}

// ID fulfills process.Sensor interface
func (p *processSensor) ID() string {
	return p.id
}

// Temperature fulfills process.Sensor interface
func (p *processSensor) Temperature() float64 {
	return p.getTemp()
}
