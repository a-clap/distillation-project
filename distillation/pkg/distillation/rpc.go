/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
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
