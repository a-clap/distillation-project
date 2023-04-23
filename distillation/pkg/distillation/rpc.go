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
