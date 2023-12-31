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
	"time"

	"embedded/pkg/ds18b20"
	"osservice"

	embeddedgpio "embedded/pkg/gpio"
	"gui/backend/ds"
	"gui/backend/gpio"
	"gui/backend/heater"
	"gui/backend/loadSaver"
	"gui/backend/parameters"
	"gui/backend/phases"
	"gui/backend/pt"

	"github.com/a-clap/logging"
)

var logger = logging.GetLogger()

type Backend struct {
	ctx          context.Context
	eventEmitter *eventEmitter
	dsChan       chan error
	ptChan       chan error
	phaseChan    chan error
	interval     time.Duration

	time       osservice.Time
	net        osservice.Net
	updater    osservice.Update
	update     Update
	waitCommit chan bool
}

func New(opts ...Option) (*Backend, error) {
	b := &Backend{
		eventEmitter: newEventEmitter(),
		interval:     time.Second,
	}

	for _, opt := range opts {
		if err := opt(b); err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Startup is called by Wails on application startup
func (b *Backend) Startup(ctx context.Context) {
	b.ctx = ctx
	b.eventEmitter.init(ctx)
	b.handleErrors()
	if err := loadSaver.Run(); err != nil {
		logger.Warn("loadSaver error ", logging.String("error", err.Error()))
	}

	heater.Refresh()
	gpio.Refresh()
	ds.Refresh()
	pt.Refresh()
	phases.Refresh()

	if errs := loadSaver.Load(); errs != nil {
		logger.Warn("Parameters Load errors", logging.Reflect("errors", errs))
	}

}

func (b *Backend) handleErrors() {
	if b.dsChan != nil {
		go func() {
			for err := range b.dsChan {
				logger.Warn("Error from DS", logging.String("error", err.Error()))
				// TODO: How to get ID based on error?
				b.eventEmitter.OnError(0)
			}
		}()
	}
	if b.ptChan != nil {
		go func() {
			for err := range b.ptChan {
				logger.Warn("Error from PT", logging.String("error", err.Error()))
				// TODO: How to get ID based on error?
				b.eventEmitter.OnError(0)
			}
		}()
	}
	if b.phaseChan != nil {
		go func() {
			for err := range b.phaseChan {
				logger.Warn("Error from Phases", logging.String("error", err.Error()))
				// TODO: How to get ID based on error?
				b.eventEmitter.OnError(0)
			}
		}()
	}
}

// Events returns Event structure - wails need to generate binding for Events methods
func (b *Backend) Events() Events {
	return Events{}
}

func (b *Backend) HeaterEnable(id string, enable bool) {
	if err := heater.Enable(id, enable); err != nil {
		logger.Error("heaterEnable", logging.String("error", err.Error()))
	}
}

func (b *Backend) HeatersGet() []parameters.Heater {
	return heater.Get()
}

func (b *Backend) DSGet() []parameters.DS {
	return ds.Get()
}

func (b *Backend) DSSetCorrection(id string, correction float64) {
	if err := ds.SetCorrection(id, correction); err != nil {
		logger.Error("SetCorrection error ", logging.String("ID", id), logging.Float64("correction", correction))
		b.eventEmitter.OnError(ErrDSSetCorrection)
	}
}

func (b *Backend) DSEnable(id string, ena bool) {
	if err := ds.Enable(id, ena); err != nil {
		logger.Error("DSEnable error ", logging.String("ID", id), logging.Bool("enable", ena))
		b.eventEmitter.OnError(ErrDSEnable)
	}
}

func (b *Backend) DSSetSamples(id string, samples uint) {
	if err := ds.SetSamples(id, samples); err != nil {
		logger.Error("DSSetSamples error ", logging.String("ID", id), logging.Uint("samples", samples))
		b.eventEmitter.OnError(ErrDSSetSamples)
	}
}

func (b *Backend) DSSetResolution(id string, res uint) {
	logger.Debug("SetResolution", logging.String("ID", id), logging.Uint("resolution", res))
	if err := ds.SetResolution(id, ds18b20.Resolution(res)); err != nil {
		logger.Error("DSSetResolution error ", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrDSSetResolution)
	}
}

func (b *Backend) DSSetName(id, name string) {
	logger.Debug("DSSetName", logging.String("ID", id), logging.String("name", name))
	if err := ds.SetName(id, name); err != nil {
		logger.Error("DSSetName error ", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrDSSetName)
	}
}

func (b *Backend) PTGet() []parameters.PT {
	return pt.Get()
}

func (b *Backend) PTSetCorrection(id string, correction float64) {
	logger.Debug("SetCorrection ", logging.String("ID", id), logging.Float64("correction", correction))
	if err := pt.SetCorrection(id, correction); err != nil {
		logger.Error("PTSetCorrection error ", logging.String("ID", id), logging.Float64("correction", correction))
		b.eventEmitter.OnError(ErrPTSetCorrection)
	}
}

func (b *Backend) PTEnable(id string, ena bool) {
	logger.Debug("PTEnable ", logging.String("ID", id), logging.Bool("enable", ena))
	if err := pt.Enable(id, ena); err != nil {
		logger.Error("PTEnable error ", logging.String("ID", id), logging.Bool("enable", ena))
		b.eventEmitter.OnError(ErrPTEnable)
	}
}

func (b *Backend) PTSetSamples(id string, samples uint) {
	logger.Debug("SetSamples ", logging.String("ID", id), logging.Uint("samples", samples))
	if err := pt.SetSamples(id, samples); err != nil {
		logger.Error("PTSetSamples error ", logging.String("ID", id), logging.Uint("samples", samples))
		b.eventEmitter.OnError(ErrPTSetSamples)
	}
}

func (b *Backend) PTSetName(id, name string) {
	logger.Debug("PTSetName", logging.String("ID", id), logging.String("name", name))
	if err := pt.SetName(id, name); err != nil {
		logger.Error("PTSetName error ", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrPTSetName)
	}
}

func (b *Backend) GPIOGet() []parameters.GPIO {
	logger.Debug("GPIOGet")
	return gpio.Get()
}

func (b *Backend) GPIOSetActiveLevel(id string, active int) {
	logger.Debug("GPIOSetActiveLevel", logging.String("id", id), logging.Int("active", active))

	if err := gpio.SetActiveLevel(id, embeddedgpio.ActiveLevel(active)); err != nil {
		logger.Error("GPIOSetActiveLevel", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrGPIOSetActiveLevel)
	}
}

func (b *Backend) GPIOSetState(id string, value bool) {
	logger.Debug("GPIOSetState", logging.String("id", id), logging.Bool("value", value))
	if err := gpio.SetState(id, value); err != nil {
		logger.Error("GPIOSetState", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrGPIOSetState)
	}
}

func (b *Backend) SaveParameters() {
	err := loadSaver.Save()
	if err != nil {
		b.eventEmitter.OnError(ErrSave)
	}
}

func (b *Backend) LoadParameters() {
	err := loadSaver.Load()
	if err != nil {
		for _, e := range err {
			logger.Debug("error on load: ", logging.String("error", e.Error()))
		}
		b.eventEmitter.OnError(ErrLoad)
	}
}
