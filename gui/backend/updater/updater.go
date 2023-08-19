package updater

import (
	"distillation/pkg/distillation"

	"gui/backend/parameters"
)

type Client interface {
	NewUpdate() (bool, string, error)

	Get() ([]distillation.PTConfig, error)
	Configure(sensor distillation.PTConfig) (distillation.PTConfig, error)
	Temperatures() ([]distillation.Temperature, error)
}

type Listener interface {
	OnPTConfigChange(parameters.PT)
	OnPTTemperatureChange(temperature parameters.Temperature)
}
