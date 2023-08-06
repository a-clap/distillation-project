/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"

	"github.com/a-clap/distillation/pkg/process"
	"embedded/pkg/embedded"
)

// ProcessPhaseCount is JSON wrapper for process.PhaseNumber
type ProcessPhaseCount struct {
	PhaseNumber uint `json:"phase_number"`
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

func (d *Distillation) Status() ProcessStatus {
	d.lastStatusMtx.Lock()
	defer d.lastStatusMtx.Unlock()
	return d.lastStatus

}
func (d *Distillation) ValidateConfig() ProcessConfigValidation {
	v := ProcessConfigValidation{Valid: true}
	err := d.Process.Validate()
	if err != nil {
		v.Valid = false
		v.Error = err.Error()
	}
	return v
}
func (d *Distillation) ConfigureProcess(cfg ProcessConfig) error {
	if !d.Process.Running() {
		// Not possible if process is not running
		if cfg.MoveToNext || cfg.Disable {
			return errors.New("process is not running")

		}
		// If user wants to enable process
		if cfg.Enable {
			s, err := d.Process.Run()
			if err != nil {
				return err
			}
			d.updateStatus(s)
			d.handleProcess()
			return nil
		}
	}

	if cfg.MoveToNext {
		s, err := d.Process.MoveToNext()
		if err != nil {
			return err
		}
		d.updateStatus(s)
		return nil
	} else if cfg.Disable {
		s, err := d.Process.Finish()
		if err != nil {
			return err
		}
		d.updateStatus(s)
		return nil
	}
	return errors.New("nothing to do")
}

func (d *Distillation) configurePhase(number uint, config ProcessPhaseConfig) error {
	// Update ios, if process is not running
	if d.Process.Running() == false {
		d.updateProcess()
	}
	return d.Process.SetPhaseConfig(number, config.PhaseConfig)

}

func (d *Distillation) updateProcess() {
	d.updateHeaters()
	d.updateOutputs()
	d.updateSensors()
}

func (d *Distillation) safeUpdateSensors() {
	if d.Process.Running() == false {
		d.updateSensors()
	}
}
func (d *Distillation) safeUpdateHeaters() {
	if d.Process.Running() == false {
		d.updateHeaters()
	}
}
func (d *Distillation) safeUpdateOutputs() {
	if d.Process.Running() == false {
		d.updateOutputs()
	}
}

func (d *Distillation) updateSensors() {
	if d.DSHandler == nil && d.PTHandler == nil {
		return
	}

	getTempDS := func(id string) func() (float64, error) {
		return func() (float64, error) {
			t := d.DSHandler.Temperature(id)
			if t.Error != 0 {
				return 0, errors.New("error on DS Temperature")
			}
			return t.Temperature, nil
		}
	}

	getTempPT := func(id string) func() (float64, error) {
		return func() (float64, error) {
			t := d.PTHandler.Temperature(id)
			if t.Error != 0 {
				return 0, errors.New("error on PT Temperature")
			}
			return t.Temperature, nil
		}
	}
	var sensors []process.Sensor
	if d.DSHandler != nil {
		for _, ds := range d.DSHandler.GetSensors() {
			if ds.Enabled {
				s := &processSensor{
					id:      ds.ID,
					getTemp: getTempDS(ds.ID),
				}
				sensors = append(sensors, s)
			}
		}
	}
	if d.PTHandler != nil {
		for _, pt := range d.PTHandler.GetSensors() {
			if pt.Enabled {
				s := &processSensor{
					id:      pt.ID,
					getTemp: getTempPT(pt.ID),
				}
				sensors = append(sensors, s)
			}
		}
	}
	d.Process.UpdateSensors(sensors)

}
func (d *Distillation) updateOutputs() {
	if d.GPIOHandler == nil {
		return
	}
	setValue := func(id string) func(value bool) error {
		return func(value bool) error {
			o, ok := d.GPIOHandler.io[id]
			if !ok {
				return ErrNoSuchID
			}
			cfg := *o
			cfg.Value = value
			_, err := d.GPIOHandler.Configure(cfg)
			return err
		}
	}

	var outputs []process.Output
	for _, out := range d.GPIOHandler.Config() {
		o := &processOutput{
			id:       out.ID,
			setValue: setValue(out.ID),
		}
		outputs = append(outputs, o)
	}
	d.Process.UpdateOutputs(outputs)

}

func (d *Distillation) updateHeaters() {
	if d.HeatersHandler == nil {
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
			_, err := d.HeatersHandler.Configure(cfg)
			return err
		}
	}

	var heaters []process.Heater
	for _, heater := range d.HeatersHandler.ConfigsGlobal() {
		if heater.Enabled {
			h := &processHeater{
				id:     heater.ID,
				setPwr: setPwr(heater.ID)}
			heaters = append(heaters, h)
		}
	}
	d.Process.UpdateHeaters(heaters)
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
	getTemp func() (float64, error)
}

// ID fulfills process.Sensor interface
func (p *processSensor) ID() string {
	return p.id
}

// Temperature fulfills process.Sensor interface
func (p *processSensor) Temperature() (float64, error) {
	return p.getTemp()
}
