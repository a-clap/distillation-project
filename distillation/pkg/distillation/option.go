/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

type Option func(*Distillation) error

func WithHeaters(heaters Heaters) Option {
	return func(h *Distillation) (err error) {
		if h.HeatersHandler, err = NewHandlerHeaters(heaters); err != nil {
			h.HeatersHandler = nil
		} else {
			h.HeatersHandler.subscribe(h.safeUpdateHeaters)
		}
		return err
	}
}
func WithGPIO(gpio GPIO) Option {
	return func(h *Distillation) (err error) {
		if h.GPIOHandler, err = NewGPIOHandler(gpio); err != nil {
			h.GPIOHandler = nil
		} else {
			h.GPIOHandler.subscribe(h.safeUpdateOutputs)
		}
		return err
	}
}

func WithDS(ds DS) Option {
	return func(h *Distillation) (err error) {
		if h.DSHandler, err = NewDSHandler(ds); err != nil {
			h.DSHandler = nil
		} else {
			h.DSHandler.subscribe(h.updateSensors)
		}
		return err
	}
}

func WithPT(pt PT) Option {
	return func(h *Distillation) (err error) {
		if h.PTHandler, err = NewPTHandler(pt); err != nil {
			h.PTHandler = nil
		} else {
			h.PTHandler.subscribe(h.updateSensors)
		}
		return err
	}
}
func WithURL(url string) Option {
	return func(d *Distillation) error {
		d.url = url
		return nil
	}
}
