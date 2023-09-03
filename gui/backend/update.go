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
	"mender"
	"osservice"

	"github.com/a-clap/logging"
)

type CheckUpdateData struct {
	NewUpdate bool     `json:"new_update"`
	Releases  []string `json:"releases"`
	ErrorCode int      `json:"error_code"`
}

type Update struct {
	NewUpdate   bool   `json:"new_update"`
	Version     string `json:"version"`
	Updating    bool   `json:"updating"`
	Downloading int    `json:"downloading"`
	Installing  int    `json:"installing"`
	Rebooting   int    `json:"rebooting"`
	Commit      bool   `json:"commit"`
}

type UpdateFinished struct {
	Version string `json:"version"`
	Status  bool   `json:"status"`
}

var _ osservice.UpdateCallbacks = (*Backend)(nil)

func (b *Backend) CheckUpdates() CheckUpdateData {
	var (
		err error
		u   CheckUpdateData
	)

	if u.NewUpdate, err = b.updater.PullReleases(); err != nil {
		u.ErrorCode = ErrUpdatePullReleases
		return u
	}

	if u.Releases, err = b.updater.AvailableReleases(); err != nil {
		u.ErrorCode = ErrAvailableReleases
		return u
	}

	return u
}

func (b *Backend) StartUpdate(name string) {
	b.update = Update{
		Updating: true,
		Version:  name,
	}

	b.waitCommit = make(chan bool)

	err := b.updater.Update(name, b)
	if err != nil {
		logger.Error("StartUpdate", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrStartUpdate)
		b.update = Update{}
		b.emitUpdate()
	}
}

func (b *Backend) Update(status mender.DeploymentStatus, progress int) {
	switch status {
	case mender.Downloading:
		b.update.Downloading = progress
	case mender.PauseBeforeInstalling:
		b.update.Downloading = 100
	case mender.Installing:
		b.update.Downloading = 100
		b.update.Installing = progress
	case mender.PauseBeforeRebooting:
		b.update.Downloading = 100
		b.update.Installing = 100
	case mender.Rebooting:
		b.update.Downloading = 100
		b.update.Installing = 100
		b.update.Rebooting = progress
	case mender.PauseBeforeCommitting:
		b.update.Downloading = 100
		b.update.Installing = 100
		b.update.Rebooting = 100
	case mender.Success, mender.Failure, mender.AlreadyInstalled:
		finished := UpdateFinished{
			Version: b.update.Version,
			Status:  status == mender.Success,
		}
		b.eventEmitter.OnUpdateFinished(finished)
		b.stopUpdate()
		return
	}

	b.emitUpdate()
}

func (b *Backend) Commit(success bool) {
	if b.update.Updating {
		b.waitCommit <- success
	}
}

func (b *Backend) Reboot(reboot bool) {
	if b.update.Updating {
		b.waitCommit <- reboot
	}
}

func (b *Backend) NextState(status mender.DeploymentStatus) bool {
	if status != mender.Success && status != mender.Rebooting {
		return true
	}

	b.update.Commit = status == mender.Success
	b.emitUpdate()
	return <-b.waitCommit
}

func (b *Backend) Error(err error) {
	logger.Error("update error", logging.String("error", err.Error()))
	b.eventEmitter.OnError(ErrUpdate)
}

func (b *Backend) StopUpdate() {
	if err := b.updater.StopUpdate(); err != nil {
		logger.Error("stop update", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrStopUpdate)
	}
	b.stopUpdate()
}

func (b *Backend) stopUpdate() {
	b.update = Update{}
	b.emitUpdate()
}

func (b *Backend) emitUpdate() {
	b.eventEmitter.OnUpdate(b.update)
}
