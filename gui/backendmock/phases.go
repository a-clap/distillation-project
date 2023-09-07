// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package backendmock

import (
	"distillation/pkg/distillation"
	"distillation/pkg/process"
	"errors"
	"fmt"

	"gui/backend/phases"
)

var _ phases.Client = (*PhasesClient)(nil)

type PhasesClient struct {
	Config  process.Config
	Process distillation.ProcessConfig
	Stats   distillation.ProcessStatus
}

// ConfigureGlobalGPIO implements phases.Client
func (p *PhasesClient) ConfigureGlobalGPIO(configs []process.GPIOConfig) ([]process.GPIOConfig, error) {
	p.Config.GlobalGPIO = configs
	return configs, nil
}

func (p *PhasesClient) Init(count uint) {
	p.Config.PhaseNumber = count
	p.Config.Sensors = []string{"s1", "s2", "s3"}
	p.Config.Phases = make([]process.PhaseConfig, count)
	for i := range p.Config.Phases {
		p.Config.Phases[i].Next = process.MoveToNextConfig{
			Type:            process.ByTime,
			SensorID:        "DS_1",
			SensorThreshold: 1.23,
			TimeLeft:        10,
		}
		p.Config.Phases[i].GPIO = make([]process.GPIOConfig, 3)
		for j := range p.Config.Phases[i].GPIO {
			p.Config.Phases[i].GPIO[j] = process.GPIOConfig{
				ID:         fmt.Sprintf("gpio %v", j),
				SensorID:   "DS_2",
				TLow:       1.23,
				THigh:      4.56,
				Hysteresis: 3.17,
				Inverted:   false,
			}
		}

		p.Config.Phases[i].Heaters = make([]process.HeaterPhaseConfig, 3)
		for j := range p.Config.Phases[i].Heaters {
			p.Config.Phases[i].Heaters[j] = process.HeaterPhaseConfig{
				ID:    fmt.Sprintf("heater %v", j),
				Power: j,
			}
		}
	}
	p.Config.GlobalGPIO = p.Config.Phases[0].GPIO
}

// ConfigurePhase implements phases.Client
func (p *PhasesClient) ConfigurePhase(phaseNumber int, setConfig distillation.ProcessPhaseConfig) (distillation.ProcessPhaseConfig, error) {
	if p.Config.PhaseNumber < uint(phaseNumber) {
		return distillation.ProcessPhaseConfig{}, errors.New("no such phase")
	}
	p.Config.Phases[phaseNumber] = setConfig.PhaseConfig
	return setConfig, nil
}

// ConfigurePhaseCount implements phases.Client
func (p *PhasesClient) ConfigurePhaseCount(count distillation.ProcessPhaseCount) (distillation.ProcessPhaseCount, error) {
	p.Init(count.PhaseNumber)
	return count, nil
}

// ConfigureProcess implements phases.Client
func (p *PhasesClient) ConfigureProcess(cfg distillation.ProcessConfig) (distillation.ProcessConfig, error) {
	p.Process = cfg
	return cfg, nil
}

// GetPhaseConfig implements phases.Client
func (p *PhasesClient) GetPhaseConfig(phaseNumber int) (distillation.ProcessPhaseConfig, error) {
	if p.Config.PhaseNumber <= uint(phaseNumber) {
		return distillation.ProcessPhaseConfig{}, errors.New("overflow")
	}
	c := distillation.ProcessPhaseConfig{PhaseConfig: p.Config.Phases[phaseNumber]}
	return c, nil
}

// GetPhaseCount implements phases.Client
func (p *PhasesClient) GetPhaseCount() (distillation.ProcessPhaseCount, error) {
	return distillation.ProcessPhaseCount{PhaseNumber: p.Config.PhaseNumber}, nil
}

// Status implements phases.Client
func (p *PhasesClient) Status() (distillation.ProcessStatus, error) {
	return p.Stats, nil
}

// ValidateConfig implements phases.Client
func (p *PhasesClient) ValidateConfig() (distillation.ProcessConfigValidation, error) {
	return distillation.ProcessConfigValidation{Valid: true}, nil
}

// GlobalConfig implements phases.Client
func (p *PhasesClient) GlobalConfig() (process.Config, error) {
	return p.Config, nil
}
