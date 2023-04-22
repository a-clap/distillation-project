/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"sync"
	"sync/atomic"
	"time"
	
	"github.com/a-clap/distillation/pkg/distillation/process"
	"github.com/a-clap/logging"
)

var (
	logger = logging.GetLogger()
)

type Distillation struct {
	HeatersHandler *HeatersHandler
	DSHandler      *DSHandler
	PTHandler      *PTHandler
	GPIOHandler    *GPIOHandler
	running        atomic.Bool
	runInterval    time.Duration
	finish         chan struct{}
	finished       chan struct{}
	Process        *process.Process
	lastStatus     ProcessStatus
	lastStatusMtx  sync.Mutex
}

func New(opts ...Option) (*Distillation, error) {
	h := &Distillation{
		running:       atomic.Bool{},
		finish:        make(chan struct{}),
		finished:      make(chan struct{}),
		runInterval:   1 * time.Second,
		lastStatusMtx: sync.Mutex{},
	}
	
	// Options
	for _, opt := range opts {
		if err := opt(h); err != nil {
			logger.Error("Option failed", logging.String("error", err.Error()))
		}
	}
	var err error
	if h.Process, err = process.New(); err != nil {
		panic(err)
	}
	
	return h, nil
}

func (h *Distillation) updateTemperatures() {
	var wg sync.WaitGroup
	if h.PTHandler != nil {
		wg.Add(1)
		go func() {
			for h.running.Load() {
				select {
				case <-h.finish:
					break
				case <-time.After(h.runInterval):
					errs := h.PTHandler.Update()
					if errs != nil {
						logger.Error("PTUpdateTemperatures failed", logging.Reflect("error", errs))
					}
				}
			}
			wg.Done()
		}()
	}
	if h.DSHandler != nil {
		wg.Add(1)
		go func() {
			for h.running.Load() {
				select {
				case <-h.finish:
					break
				case <-time.After(h.runInterval):
					errs := h.DSHandler.Update()
					if errs != nil {
						logger.Error("DSUpdateTemperatures failed", logging.Reflect("error", errs))
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(h.finish)
}

func (h *Distillation) handleProcess() {
	go func() {
		for h.Process.Running() {
			select {
			case <-time.After(h.runInterval):
				s, err := h.Process.Process()
				if err != nil {
					logger.Error("HandleProcess", logging.String("error", err.Error()))
				} else {
					h.updateStatus(s)
				}
			}
		}
	}()
	
}

func (h *Distillation) updateStatus(s process.Status) {
	h.lastStatusMtx.Lock()
	h.lastStatus.Status = s
	h.lastStatusMtx.Unlock()
}
