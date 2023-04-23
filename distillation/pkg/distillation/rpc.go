/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"time"
	
	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
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
