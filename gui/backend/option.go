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
	"osservice"

	"gui/backend/ds"
	"gui/backend/gpio"
	"gui/backend/heater"
	"gui/backend/loadSaver"
	"gui/backend/phases"
	"gui/backend/pt"
	"gui/backend/wifi"
)

type Option func(b *Backend) error

func WithHeaterClient(c heater.Client) Option {
	return func(b *Backend) error {
		return heater.Init(c)
	}
}

func WithDSClient(c ds.Client) Option {
	return func(b *Backend) error {
		b.dsChan = make(chan error, 10)
		if err := ds.Init(c, b.dsChan, b.interval); err != nil {
			return err
		}
		ds.Run()
		return nil
	}
}

func WithPTClient(c pt.Client) Option {
	return func(b *Backend) error {
		b.ptChan = make(chan error, 10)
		if err := pt.Init(c, b.ptChan, b.interval); err != nil {
			return err
		}
		pt.Run()
		return nil
	}
}

func WithGPIOClient(c gpio.Client) Option {
	return func(b *Backend) error {
		return gpio.Init(c)
	}
}

func WithWifi(c wifi.Client) Option {
	return func(b *Backend) error {
		wifi.Init(c)
		return nil
	}
}

func WithPhaseClient(c phases.Client) Option {
	return func(b *Backend) error {
		b.phaseChan = make(chan error, 10)
		phases.Init(c, b.phaseChan, b.interval)
		phases.Run()
		return nil
	}
}

func WithLoadSaver(c loadSaver.LoadSaver) Option {
	return func(b *Backend) error {
		loadSaver.Init(c)
		return nil
	}
}

func WithTime(ts osservice.Time) Option {
	return func(b *Backend) error {
		b.time = ts
		return nil
	}
}

func WithNet(net osservice.Net) Option {
	return func(b *Backend) error {
		b.net = net
		return nil
	}
}

func WithUpdate(update osservice.Update) Option {
	return func(b *Backend) error {
		b.updater = update
		return nil
	}
}
