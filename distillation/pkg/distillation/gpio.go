/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"
	
	"github.com/a-clap/embedded/pkg/embedded"
)

var (
	ErrNoGPIOInterface = errors.New("no GPIO interface")
)

type GPIO interface {
	Get() ([]embedded.GPIOConfig, error)
	Configure(c embedded.GPIOConfig) (embedded.GPIOConfig, error)
}

type GPIOError struct {
	ID  string `json:"ID"`
	Op  string `json:"op"`
	Err string `json:"error"`
}

func (g *GPIOError) Error() string {
	if g.Err == "" {
		return "<nil>"
	}
	s := g.Op
	if g.ID != "" {
		s += ":" + g.ID
	}
	s += ": " + g.Err
	return s
}

type GPIOConfig struct {
	embedded.GPIOConfig
}

type GPIOHandler struct {
	GPIO    GPIO
	io      map[string]*GPIOConfig
	clients []gpioCallback
}

type gpioCallback func()

func NewGPIOHandler(io GPIO) (*GPIOHandler, error) {
	g := &GPIOHandler{
		GPIO:    io,
		io:      make(map[string]*GPIOConfig, 0),
		clients: make([]gpioCallback, 0),
	}
	if err := g.init(); err != nil {
		return nil, err
	}
	return g, nil
}

func (g *GPIOHandler) Config() []GPIOConfig {
	configs := make([]GPIOConfig, 0, len(g.io))
	for _, value := range g.io {
		configs = append(configs, *value)
	}
	return configs
}

func (g *GPIOHandler) Configure(cfg GPIOConfig) (GPIOConfig, error) {
	c := GPIOConfig{}
	io, ok := g.io[cfg.ID]
	if !ok {
		return c, &GPIOError{ID: cfg.ID, Op: "Configure", Err: ErrNoSuchID.Error()}
	}
	newcfg, err := g.GPIO.Configure(cfg.GPIOConfig)
	if err != nil {
		return c, &GPIOError{ID: cfg.ID, Op: "Configure", Err: err.Error()}
	}
	io.GPIOConfig = newcfg
	g.notify()
	return *io, nil
	
}
func (g *GPIOHandler) init() error {
	if g.GPIO == nil {
		return &GPIOError{Op: "init", Err: ErrNoGPIOInterface.Error()}
	}
	
	ios, err := g.GPIO.Get()
	if err != nil {
		return &GPIOError{Op: "init.Get", Err: err.Error()}
	}
	
	for _, io := range ios {
		g.io[io.ID] = &GPIOConfig{GPIOConfig: io}
	}
	return nil
}

func (g *GPIOHandler) subscribe(cb gpioCallback) {
	g.clients = append(g.clients, cb)
}
func (g *GPIOHandler) notify() {
	for _, client := range g.clients {
		client()
	}
}
