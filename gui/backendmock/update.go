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

package backendmock

import (
	"fmt"
	"mender"
	"osservice"
	"sync/atomic"
	"time"
)

type Update struct {
	newUpdate bool
	updating  atomic.Bool
	stop      chan struct{}
	callbacks osservice.UpdateCallbacks
}

var (
	_         osservice.Update = (*Update)(nil)
	releases                   = []string{"release1"}
	stepDelay                  = 10 * time.Millisecond
)

func (u *Update) ContinueUpdate() (bool, string) {
	return false, ""
}

func (u *Update) PullReleases() (bool, error) {
	u.newUpdate = !u.newUpdate
	return u.newUpdate, nil
}

func (u *Update) AvailableReleases() ([]string, error) {
	time.Sleep(1 * time.Second)
	return releases, nil
}

func (u *Update) Update(artifactName string, callbacks osservice.UpdateCallbacks) error {
	if u.updating.Load() {
		return fmt.Errorf("already updating")
	}

	u.updating.Store(true)
	u.stop = make(chan struct{})
	u.callbacks = callbacks

	go u.update()

	return nil
}

func (u *Update) StopUpdate() error {
	if u.updating.Load() {
		u.updating.Store(false)
		close(u.stop)
	}
	return nil
}

func (u *Update) update() {
	updateState := func(state mender.DeploymentStatus, maxProgress int) {
		progress := 0
		for u.updating.Load() && progress < maxProgress {
			select {
			case <-u.stop:
				u.updating.Store(false)
			case <-time.After(stepDelay):
				progress++
				u.callbacks.Update(state, progress)
			}
		}
	}

	for _, state := range []mender.DeploymentStatus{mender.Downloading, mender.Installing, mender.Rebooting} {
		updateState(state, 100)
	}

	u.updating.Store(false)
}
