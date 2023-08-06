/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"sync"
	"sync/atomic"
	"time"

	"distillation/pkg/process"
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
	h.Process = process.New()

	// Update process with enabled components
	h.safeUpdateHeaters()
	h.safeUpdateSensors()
	h.safeUpdateOutputs()

	return h, nil
}

func (d *Distillation) Run() {
	d.running.Store(true)
	go d.updateTemperatures()
}

func (d *Distillation) Close() {
	d.running.Store(false)
	close(d.finish)
	for range d.finished {
	}
}

func (d *Distillation) updateTemperatures() {
	var wg sync.WaitGroup
	if d.PTHandler != nil {
		wg.Add(1)
		go func() {
			for d.running.Load() {
				select {
				case <-d.finish:
					break
				case <-time.After(d.runInterval):
					errs := d.PTHandler.Update()
					if errs != nil {
						logger.Error("PTUpdateTemperatures failed", logging.Reflect("error", errs))
					}
				}
			}
			wg.Done()
		}()
	}
	if d.DSHandler != nil {
		wg.Add(1)
		go func() {
			for d.running.Load() {
				select {
				case <-d.finish:
					break
				case <-time.After(d.runInterval):
					errs := d.DSHandler.Update()
					if errs != nil {
						logger.Error("DSUpdateTemperatures failed", logging.Reflect("error", errs))
					}
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	close(d.finished)
}

func (d *Distillation) handleProcess() {
	go func() {
		for d.Process.Running() {
			select {
			case <-time.After(d.runInterval):
				s, err := d.Process.Process()
				if err != nil {
					logger.Error("HandleProcess", logging.String("error", err.Error()))
				} else {
					d.updateStatus(s)
				}
			}
		}
	}()

}

func (d *Distillation) updateStatus(s process.Status) {
	d.lastStatusMtx.Lock()
	d.lastStatus.Status = s
	d.lastStatusMtx.Unlock()
}
