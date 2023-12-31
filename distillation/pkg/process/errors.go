/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package process

import (
	"errors"
)

var (
	ErrNoHeaters                         = errors.New("can't execute process without heaters")
	ErrNoTSensors                        = errors.New("can't execute process without temperature sensors")
	ErrNoSuchPhase                       = errors.New("requested phase number doesn't exist")
	ErrWrongHeaterID                     = errors.New("requested heater ID doesn't exist")
	ErrWrongSensorID                     = errors.New("requested sensor ID doesn't exist")
	ErrWrongGpioID                       = errors.New("requested gpio ID doesn't exist")
	ErrWrongHeaterPower                  = errors.New("power must be in range <0, 100>")
	ErrDifferentGPIOSConfig              = errors.New("different number of gpio configs and gpios")
	ErrHeaterConfigDiffersFromHeatersLen = errors.New("different number of heaters configs and heaters")
	ErrByTimeWrongTime                   = errors.New("timeleft must be greater than 0")
	ErrByTemperatureWrongID              = errors.New("chosen ByTemperature, but id doesn't exist")
	ErrUnknownType                       = errors.New("MoveToNextType unknown")
	ErrNotRunning                        = errors.New("not running")
	ErrTooManyErrorOnTemperature         = errors.New("too many consecutive errors on Temperature")
	ErrSetPower                          = errors.New("on heater SetPower")
	ErrSetValue                          = errors.New("on output Set")
	ErrDuplicatedID                      = errors.New("duplicated ID in config")
)
