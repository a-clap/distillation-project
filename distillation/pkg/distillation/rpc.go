/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"time"

	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
	"github.com/a-clap/distillation/pkg/process"
	"github.com/a-clap/embedded/pkg/ds18b20"
	"github.com/a-clap/embedded/pkg/embedded"
	"github.com/a-clap/embedded/pkg/gpio"
	"github.com/a-clap/embedded/pkg/max31865"
)

func gpioConfigToRPC(config *GPIOConfig) *distillationproto.GPIOConfig {
	return &distillationproto.GPIOConfig{
		ID:          config.ID,
		Direction:   int32(config.Direction),
		ActiveLevel: int32(config.ActiveLevel),
		Value:       config.Value,
	}
}

func rpcToGPIOConfig(config *distillationproto.GPIOConfig) GPIOConfig {
	return GPIOConfig{GPIOConfig: embedded.GPIOConfig{
		Config: gpio.Config{
			ID:          config.ID,
			Direction:   gpio.Direction(config.Direction),
			ActiveLevel: gpio.ActiveLevel(config.ActiveLevel),
			Value:       config.Value,
		},
	}}
}

func rpcToDSConfig(elem *distillationproto.DSConfig) DSConfig {
	return DSConfig{
		DSSensorConfig: embedded.DSSensorConfig{
			Enabled: elem.Enabled,
			SensorConfig: ds18b20.SensorConfig{
				Name:         elem.Name,
				ID:           elem.ID,
				Correction:   float64(elem.Correction),
				Resolution:   ds18b20.Resolution(elem.Resolution),
				PollInterval: time.Duration(elem.PollInterval),
				Samples:      uint(elem.Samples),
			},
		},
		temps: embedded.DSTemperature{},
	}
}

func dsConfigToRPC(d *DSConfig) *distillationproto.DSConfig {
	return &distillationproto.DSConfig{
		ID:           d.ID,
		Name:         d.Name,
		Correction:   float32(d.Correction),
		Resolution:   int32(d.Resolution),
		PollInterval: int32(d.PollInterval),
		Samples:      uint32(d.Samples),
		Enabled:      d.Enabled,
	}
}

func rpcToDSTemperature(r *distillationproto.DSTemperatures) []DSTemperature {
	temperatures := make([]DSTemperature, len(r.Temps))
	for i, temp := range r.Temps {
		temperatures[i] = DSTemperature{
			ID:          temp.ID,
			Temperature: float64(temp.Temperature),
		}
	}
	return temperatures
}

func dsTemperatureToRPC(t []DSTemperature) *distillationproto.DSTemperatures {
	temperatures := &distillationproto.DSTemperatures{}
	temperatures.Temps = make([]*distillationproto.DSTemperature, len(t))
	for i, temp := range t {
		temperatures.Temps[i] = &distillationproto.DSTemperature{
			ID:          temp.ID,
			Temperature: float32(temp.Temperature),
		}
	}
	return temperatures
}

func heaterConfigToRPC(config *HeaterConfigGlobal) *distillationproto.HeaterConfig {
	return &distillationproto.HeaterConfig{
		ID:      config.ID,
		Enabled: config.Enabled,
	}
}

func rpcToHeaterConfig(config *distillationproto.HeaterConfig) HeaterConfigGlobal {
	return HeaterConfigGlobal{
		ID:      config.ID,
		Enabled: config.Enabled,
	}
}

func rpcToPTConfig(elem *distillationproto.PTConfig) PTConfig {
	return PTConfig{
		PTSensorConfig: embedded.PTSensorConfig{
			Enabled: elem.Enabled,
			SensorConfig: max31865.SensorConfig{
				Name:         elem.Name,
				ID:           elem.ID,
				Correction:   float64(elem.Correction),
				ASyncPoll:    elem.Async,
				PollInterval: time.Duration(elem.PollInterval),
				Samples:      uint(elem.Samples),
			},
		},
		temps: embedded.PTTemperature{},
	}
}

func ptConfigToRPC(d *PTConfig) *distillationproto.PTConfig {
	return &distillationproto.PTConfig{
		ID:           d.ID,
		Name:         d.Name,
		Correction:   float32(d.Correction),
		PollInterval: int32(d.PollInterval),
		Samples:      uint32(d.Samples),
		Enabled:      d.Enabled,
		Async:        d.ASyncPoll,
	}
}

func rpcToPTTemperature(r *distillationproto.PTTemperatures) []PTTemperature {
	temperatures := make([]PTTemperature, len(r.Temps))
	for i, temp := range r.Temps {
		temperatures[i] = PTTemperature{
			ID:          temp.ID,
			Temperature: float64(temp.Temperature),
		}
	}
	return temperatures
}

func ptTemperatureToRPC(t []PTTemperature) *distillationproto.PTTemperatures {
	temperatures := &distillationproto.PTTemperatures{}
	temperatures.Temps = make([]*distillationproto.PTTemperature, len(t))
	for i, temp := range t {
		temperatures.Temps[i] = &distillationproto.PTTemperature{
			ID:          temp.ID,
			Temperature: float32(temp.Temperature),
		}
	}
	return temperatures
}

func rpcToProcessPhaseCount(cnt *distillationproto.ProcessPhaseCount) ProcessPhaseCount {
	return ProcessPhaseCount{PhaseNumber: uint(cnt.Count)}
}

func rpcToProcessPhaseConfig(cfg *distillationproto.ProcessPhaseConfig) ProcessPhaseConfig {
	c := ProcessPhaseConfig{PhaseConfig: process.PhaseConfig{
		Next: process.MoveToNextConfig{
			Type:            process.MoveToNextType(cfg.Next.Type),
			SensorID:        cfg.Next.SensorID,
			SensorThreshold: float64(cfg.Next.SensorThreshold),
			TimeLeft:        cfg.Next.TimeLeft,
		},
		Heaters: make([]process.HeaterPhaseConfig, len(cfg.Heaters)),
		GPIO:    make([]process.GPIOConfig, len(cfg.GPIO)),
	},
	}
	for i, heater := range cfg.Heaters {
		c.Heaters[i] = process.HeaterPhaseConfig{
			ID:    heater.ID,
			Power: int(heater.Power),
		}
	}
	for i, gp := range cfg.GPIO {
		c.GPIO[i] = rpcToGPIOPhaseConfig(gp)
	}
	return c
}

func rpcToGPIOPhaseConfig(gp *distillationproto.GPIOPhaseConfig) process.GPIOConfig {
	return process.GPIOConfig{
		Enabled:    gp.Enabled,
		ID:         gp.ID,
		SensorID:   gp.SensorID,
		TLow:       float64(gp.TLow),
		THigh:      float64(gp.THigh),
		Hysteresis: float64(gp.Hysteresis),
		Inverted:   gp.Inverted,
	}
}
func gpioPhaseConfigToRpc(gp process.GPIOConfig) *distillationproto.GPIOPhaseConfig {
	return &distillationproto.GPIOPhaseConfig{
		ID:         gp.ID,
		SensorID:   gp.SensorID,
		TLow:       float32(gp.TLow),
		THigh:      float32(gp.THigh),
		Hysteresis: float32(gp.Hysteresis),
		Inverted:   gp.Inverted,
		Enabled:    gp.Enabled,
	}
}
func processPhaseConfigToRpc(number int, config ProcessPhaseConfig) *distillationproto.ProcessPhaseConfig {
	cfg := &distillationproto.ProcessPhaseConfig{
		Number: &distillationproto.PhaseNumber{Number: int32(number)},
		Next: &distillationproto.MoveToNextConfig{
			Type:            int32(config.Next.Type),
			SensorID:        config.Next.SensorID,
			SensorThreshold: float32(config.Next.SensorThreshold),
			TimeLeft:        config.Next.TimeLeft,
		},
		Heaters: make([]*distillationproto.HeaterPhaseConfig, len(config.Heaters)),
		GPIO:    make([]*distillationproto.GPIOPhaseConfig, len(config.GPIO)),
	}
	for i, heater := range config.Heaters {
		cfg.Heaters[i] = &distillationproto.HeaterPhaseConfig{
			ID:    heater.ID,
			Power: int32(heater.Power),
		}
	}
	for i, gp := range config.GPIO {
		cfg.GPIO[i] = gpioPhaseConfigToRpc(gp)
	}

	return cfg
}

func processConfigToRpc(cfg ProcessConfig) *distillationproto.ProcessConfig {
	return &distillationproto.ProcessConfig{
		Enable:     cfg.Enable,
		MoveToNext: cfg.MoveToNext,
		Disable:    cfg.Disable,
	}
}

func rpcToProcessConfig(conf *distillationproto.ProcessConfig) ProcessConfig {
	return ProcessConfig{
		Enable:     conf.Enable,
		MoveToNext: conf.MoveToNext,
		Disable:    conf.Disable,
	}
}

func rpcToProcessStatus(status *distillationproto.ProcessStatus) ProcessStatus {
	s := ProcessStatus{process.Status{
		Running:     status.Running,
		Done:        status.Done,
		PhaseNumber: uint(status.PhaseNumber),
		StartTime:   time.Unix(status.StartTime, 0),
		EndTime:     time.Unix(status.EndTime, 0),
		Next: process.MoveToNextStatus{
			Type:     process.MoveToNextType(status.Next.Type),
			TimeLeft: status.Next.TimeLeft,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        status.Next.Temperature.SensorID,
				SensorThreshold: float64(status.Next.Temperature.SensorThreshold),
			},
		},
		Heaters:     make([]process.HeaterPhaseStatus, len(status.Heaters)),
		Temperature: make([]process.TemperaturePhaseStatus, len(status.Heaters)),
		GPIO:        make([]process.GPIOPhaseStatus, len(status.Heaters)),
		Errors:      make([]string, len(status.Heaters)),
	}}
	for i, heater := range status.Heaters {
		s.Heaters[i] = process.HeaterPhaseStatus{HeaterPhaseConfig: process.HeaterPhaseConfig{
			ID:    heater.ID,
			Power: int(heater.Power),
		}}
	}
	for i, t := range status.Temperature {
		s.Temperature[i] = process.TemperaturePhaseStatus{
			ID:          t.ID,
			Temperature: float64(t.Temperature),
		}
	}

	for i, io := range status.GPIO {
		s.GPIO[i] = process.GPIOPhaseStatus{
			ID:    io.ID,
			State: io.State,
		}
	}
	for i, err := range status.Errors {
		s.Errors[i] = err
	}
	return s
}

func processStatusToRPC(status ProcessStatus) *distillationproto.ProcessStatus {
	s := &distillationproto.ProcessStatus{
		Running:     status.Running,
		Done:        status.Done,
		PhaseNumber: int32(status.PhaseNumber),
		StartTime:   status.StartTime.Unix(),
		EndTime:     status.StartTime.Unix(),
		Next: &distillationproto.MoveToNextStatus{
			Type:     int32(status.Next.Type),
			TimeLeft: status.Next.TimeLeft,
			Temperature: &distillationproto.MoveToNextStatusTemperature{
				SensorID:        status.Next.Temperature.SensorID,
				SensorThreshold: float32(status.Next.Temperature.SensorThreshold),
			}},
		Heaters:     make([]*distillationproto.HeaterPhaseStatus, len(status.Heaters)),
		Temperature: make([]*distillationproto.TemperaturePhaseStatus, len(status.Temperature)),
		GPIO:        make([]*distillationproto.GPIOPhaseStatus, len(status.GPIO)),
		Errors:      make([]string, len(status.Heaters)),
	}
	for i, heater := range status.Heaters {
		s.Heaters[i] = &distillationproto.HeaterPhaseStatus{
			ID:    heater.ID,
			Power: int32(heater.Power),
		}
	}

	for i, t := range status.Temperature {
		s.Temperature[i] = &distillationproto.TemperaturePhaseStatus{
			ID:          t.ID,
			Temperature: float32(t.Temperature),
		}
	}

	for i, io := range status.GPIO {
		s.GPIO[i] = &distillationproto.GPIOPhaseStatus{
			ID:    io.ID,
			State: io.State,
		}
	}

	for i, err := range status.Errors {
		s.Errors[i] = err
	}
	return s

}
