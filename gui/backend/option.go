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
