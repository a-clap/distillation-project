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

package backend

import (
	"context"

	"distillation/pkg/distillation"
	"distillation/pkg/process"

	"gui/backend/ds"
	"gui/backend/gpio"
	"gui/backend/heater"
	"gui/backend/parameters"
	"gui/backend/phases"
	"gui/backend/pt"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type ProcessStatus struct {
	StartTime int64 `json:"unix_start_time"`
	EndTime   int64 `json:"unix_end_time"`
	distillation.ProcessStatus
}

type eventEmitter struct {
	ctx context.Context
}

func newEventEmitter() *eventEmitter {
	return &eventEmitter{}
}

func (e *eventEmitter) init(ctx context.Context) {
	e.ctx = ctx
	heater.AddListener(e)
	ds.AddListener(e)
	pt.AddListener(e)
	gpio.AddListener(e)
	phases.AddListener(e)
}

// OnHeaterChange implements heater.Listener
func (e *eventEmitter) OnHeaterChange(h parameters.Heater) {
	runtime.EventsEmit(e.ctx, NotifyHeaters, h)
}

// OnDSConfigChange implements ds.Listener
func (e *eventEmitter) OnDSConfigChange(c parameters.DS) {
	runtime.EventsEmit(e.ctx, NotifyDSConfig, c)
}

// OnDSTemperatureChange implements ds.Listener
func (e *eventEmitter) OnDSTemperatureChange(t parameters.Temperature) {
	runtime.EventsEmit(e.ctx, NotifyDSTemperature, t)
}

// OnPTConfigChange implements pt.Listener
func (e *eventEmitter) OnPTConfigChange(p parameters.PT) {
	runtime.EventsEmit(e.ctx, NotifyPTConfig, p)
}

// OnPTTemperatureChange implements pt.Listener
func (e *eventEmitter) OnPTTemperatureChange(t parameters.Temperature) {
	runtime.EventsEmit(e.ctx, NotifyPTTemperature, t)
}

// OnGPIOChange implements gpio.Listener
func (e *eventEmitter) OnGPIOChange(config parameters.GPIO) {
	runtime.EventsEmit(e.ctx, NotifyGPIO, config)
}

func (e *eventEmitter) OnError(errID int) {
	runtime.EventsEmit(e.ctx, NotifyError, errID)
}

// OnConfigChange implements phases.Listener
func (e *eventEmitter) OnConfigChange(c distillation.ProcessConfig) {
	runtime.EventsEmit(e.ctx, NotifyPhasesConfig, c)
}

// OnConfigValidate implements phases.Listener
func (e *eventEmitter) OnConfigValidate(validation distillation.ProcessConfigValidation) {
	runtime.EventsEmit(e.ctx, NotifyPhasesValidate, validation)
}

// OnPhaseConfigChange implements phases.Listener
func (e *eventEmitter) OnPhaseConfigChange(phaseNumber int, cfg distillation.ProcessPhaseConfig) {
	runtime.EventsEmit(e.ctx, NotifyPhasesPhaseConfig, phaseNumber, cfg)
}

// OnPhasesCountChange implements phases.Listener
func (e *eventEmitter) OnPhasesCountChange(count distillation.ProcessPhaseCount) {
	runtime.EventsEmit(e.ctx, NotifyPhasesPhaseCount, count)
}

// OnStatusChange implements phases.Listener
func (e *eventEmitter) OnStatusChange(status distillation.ProcessStatus) {
	p := ProcessStatus{
		StartTime:     status.StartTime.Unix(),
		EndTime:       status.EndTime.Unix(),
		ProcessStatus: status,
	}

	runtime.EventsEmit(e.ctx, NotifyPhasesStatus, p)
}

// OnGlobalConfig implements phases.Listener
func (e *eventEmitter) OnGlobalConfig(c process.Config) {
	runtime.EventsEmit(e.ctx, NotifyGlobalConfig, c)
}

func (e *eventEmitter) OnUpdate(u Update) {
	runtime.EventsEmit(e.ctx, NotifyUpdate, u)
}

func (e *eventEmitter) OnUpdateStatus(u UpdateStateStatus) {
	runtime.EventsEmit(e.ctx, NotifyUpdateStatus, u)
}

func (e *eventEmitter) UpdateNextState(u UpdateNextState) {
	runtime.EventsEmit(e.ctx, NotifyUpdateNextState, u)
}
