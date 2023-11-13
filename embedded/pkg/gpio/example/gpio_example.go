/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"embedded/pkg/gpio"
)

func heartbeat(pin gpio.Pin, closeCh chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	out, err := gpio.Output(pin, "", false)
	if err != nil {
		panic(err)
	}

	states := []struct {
		delay time.Duration
		value bool
	}{
		{
			delay: 120 * time.Millisecond,
			value: true,
		},
		{
			delay: 60 * time.Millisecond,
			value: false,
		},
		{
			delay: 160 * time.Millisecond,
			value: true,
		},
		{
			delay: 300 * time.Millisecond,
			value: false,
		},
	}

	running := atomic.Bool{}
	running.Store(true)
	for running.Load() {

		for _, state := range states {
			err = out.Set(state.value)
			if err != nil {
				log.Println(err)
			}

			select {
			case <-time.After(state.delay):
			case <-closeCh:
				running.Store(false)

			}

			if !running.Load() {
				break
			}
		}
	}
	wg.Done()

}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	wg := sync.WaitGroup{}

	pwr := gpio.GetBananaPin(gpio.PWR_LED)
	ch := make(chan struct{})

	for _, pin := range []gpio.Pin{
		gpio.GetBananaPin(gpio.CON2_P40),
		gpio.GetBananaPin(gpio.CON2_P22),
		gpio.GetBananaPin(gpio.CON2_P38),
		pwr,
	} {
		go heartbeat(pin, ch, &wg)
	}

	log.Println("Running, waiting for CTRL+C")

	<-sig
	close(ch)

	wg.Wait()

	log.Println("Finished")
}
