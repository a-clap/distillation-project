/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package phases

import (
	"distillation/pkg/distillation"
	"distillation/pkg/process"

	"golang.org/x/exp/slices"
)

// notifyProcessCount notifies listeners about distillation.ProcessPhaseCount change
func notifyProcessCount(count distillation.ProcessPhaseCount) {
	for _, listener := range handler.listeners {
		listener.OnPhasesCountChange(count)
	}
}

func notifyGlobalConfig(c process.Config) {
	for _, listener := range handler.listeners {
		listener.OnGlobalConfig(c)
	}
}

func notifyConfigChange(phaseNumber int, cfg distillation.ProcessPhaseConfig) {
	for _, listener := range handler.listeners {
		listener.OnPhaseConfigChange(phaseNumber, cfg)
	}
}

func notifyValidate(v distillation.ProcessConfigValidation) {
	for _, listener := range handler.listeners {
		listener.OnConfigValidate(v)
	}
}

func notifyStatus(s distillation.ProcessStatus) {
	// Sort sensors
	slices.SortStableFunc(s.Temperature, func(i, j process.TemperaturePhaseStatus) int {
		if i.ID > j.ID {
			return 1
		}
		if i.ID < j.ID {
			return -1
		}
		return 0
	})

	// Sort GPIO
	slices.SortStableFunc(s.GPIO, func(i, j process.GPIOPhaseStatus) int {
		if i.ID > j.ID {
			return 1
		}
		if i.ID < j.ID {
			return -1
		}
		return 0

	})
	// And heaters
	slices.SortStableFunc(s.Heaters, func(i, j process.HeaterPhaseStatus) int {
		if i.ID > j.ID {
			return 1
		}
		if i.ID < j.ID {
			return -1
		}
		return 0
	})
	for _, listener := range handler.listeners {
		listener.OnStatusChange(s)
	}
}

func notifyConfig(c distillation.ProcessConfig) {
	for _, listener := range handler.listeners {
		listener.OnConfigChange(c)
	}
}

// notifyError notifies user about error, which can't be return via return
func notifyError(e error) {
	if handler.err != nil {
		// error is discarded, if channel is nil
		handler.err <- e
	}
}
