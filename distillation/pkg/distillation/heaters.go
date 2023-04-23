/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"
	
	"github.com/a-clap/embedded/pkg/embedded"
)

type Heaters interface {
	Get() ([]embedded.HeaterConfig, error)
	Configure(heater embedded.HeaterConfig) (embedded.HeaterConfig, error)
}

type HeaterConfigGlobal struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
}

type HeaterConfig struct {
	global HeaterConfigGlobal
	embedded.HeaterConfig
}

type HeatersHandler struct {
	Heaters Heaters
	heaters map[string]*HeaterConfig
	clients []heaterCallback
}

type HeaterError struct {
	ID  string `json:"ID"`
	Op  string `json:"op"`
	Err string `json:"error"`
}

type heaterCallback func()

func (e *HeaterError) Error() string {
	if e.Err == "" {
		return "<nil>"
	}
	s := e.Op
	if e.ID != "" {
		s += ":" + e.ID
	}
	s += ": " + e.Err
	return s
}

func NewHandlerHeaters(heaters Heaters) (*HeatersHandler, error) {
	h := &HeatersHandler{
		Heaters: heaters,
		heaters: make(map[string]*HeaterConfig),
		clients: make([]heaterCallback, 0),
	}
	if err := h.init(); err != nil {
		return nil, &HeaterError{ID: "", Op: "NewHandlerHeaters.init", Err: err.Error()}
	}
	return h, nil
}

func (h *HeatersHandler) init() error {
	if h.Heaters == nil {
		return errors.New("no Heater interface")
	}
	heaters, err := h.Heaters.Get()
	if err != nil {
		return err
	}
	
	for _, heater := range heaters {
		id := heater.ID
		cfg := HeaterConfig{
			global: HeaterConfigGlobal{
				ID:      id,
				Enabled: false,
			},
			HeaterConfig: embedded.HeaterConfig{
				ID:      id,
				Enabled: false,
				Power:   0,
			},
		}
		
		h.heaters[id] = &cfg
		if _, err = h.Configure(cfg); err != nil {
			return err
		}
		
	}
	return nil
}

func (h *HeatersHandler) ConfigsGlobal() []HeaterConfigGlobal {
	heaters := make([]HeaterConfigGlobal, 0, len(h.heaters))
	for _, v := range h.heaters {
		heaters = append(heaters, v.global)
	}
	return heaters
}

func (h *HeatersHandler) ConfigureGlobal(cfg HeaterConfigGlobal) (HeaterConfigGlobal, error) {
	c := HeaterConfigGlobal{}
	maybeHeater, ok := h.heaters[cfg.ID]
	if !ok {
		return c, &HeaterError{ID: cfg.ID, Op: "ConfigureGlobal", Err: ErrNoSuchID.Error()}
	}
	
	if maybeHeater.global.Enabled != cfg.Enabled {
		maybeHeater.global.Enabled = cfg.Enabled
		// Do we need to disable heater?
		if !maybeHeater.global.Enabled && maybeHeater.HeaterConfig.Enabled {
			maybeHeater.HeaterConfig.Enabled = false
			if _, err := h.Configure(*maybeHeater); err != nil {
				return c, err
			}
		}
	}
	h.notify()
	return maybeHeater.global, nil
}

func (h *HeatersHandler) Configs() []HeaterConfig {
	heaters := make([]HeaterConfig, 0, len(h.heaters))
	for _, v := range h.heaters {
		if v.global.Enabled {
			heaters = append(heaters, *v)
		}
	}
	return heaters
}

func (h *HeatersHandler) Config(id string) (embedded.HeaterConfig, error) {
	cfg, ok := h.heaters[id]
	if !ok {
		return embedded.HeaterConfig{}, &HeaterError{ID: id, Op: "Config", Err: ErrNoSuchID.Error()}
	}
	return cfg.HeaterConfig, nil
}

func (h *HeatersHandler) ConfigGlobal(id string) (HeaterConfigGlobal, error) {
	cfg, ok := h.heaters[id]
	if !ok {
		return HeaterConfigGlobal{}, &HeaterError{ID: id, Op: "Config", Err: ErrNoSuchID.Error()}
	}
	return cfg.global, nil
}

func (h *HeatersHandler) Configure(cfg HeaterConfig) (HeaterConfig, error) {
	c := HeaterConfig{}
	maybeHeater, ok := h.heaters[cfg.ID]
	if !ok {
		return c, &HeaterError{ID: cfg.ID, Op: "ConfigureGlobal", Err: ErrNoSuchID.Error()}
	}
	// Global has to be set
	maybeHeater.Enabled = maybeHeater.global.Enabled && cfg.Enabled
	maybeHeater.Power = cfg.Power
	newConfig, err := h.Heaters.Configure(maybeHeater.HeaterConfig)
	if err != nil {
		maybeHeater.HeaterConfig = newConfig
	}
	
	return *maybeHeater, err
}

func (h *HeatersHandler) subscribe(cb heaterCallback) {
	h.clients = append(h.clients, cb)
}
func (h *HeatersHandler) notify() {
	for _, client := range h.clients {
		client()
	}
}
