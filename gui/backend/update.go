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

type UpdateData struct {
	Releases  []string `json:"releases"`
	ErrorCode int      `json:"error_code"`
}

type UpdateNextState struct {
	State mender.DeploymentStatus `json:"state"`
}

type UpdateStateStatus struct {
	State    mender.DeploymentStatus `json:"state"`
	Progress int                     `json:"progress"`
}

type Update struct {
	Updating bool   `json:"updating"`
	Release  string `json:"release"`
	Success  bool   `json:"success"`
}

var _ osservice.UpdateCallbacks = (*Backend)(nil)

func (b *Backend) ContinueUpdate() {
	logger.Debug("try: ContinueUpdate")
	if finish, release := b.updater.ContinueUpdate(); finish {
		logger.Debug("Continuing", logging.String("release", release))
		b.StartUpdate(release)
	}
}

func (b *Backend) CheckUpdates() UpdateData {
	var (
		err error
		u   UpdateData
	)

	if _, err = b.updater.PullReleases(); err != nil {
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

	logger.Debug("StartUpdate", logging.String("name", name))
	b.update = Update{
		Updating: true,
		Release:  name,
	}

	b.waitCommit = make(chan bool)

	err := b.updater.Update(name, b)
	if err != nil {
		logger.Error("StartUpdate", logging.String("error", err.Error()))
		b.eventEmitter.OnError(ErrStartUpdate)
		b.finishUpdate(false)
	} else {
		b.eventEmitter.OnUpdate(b.update)
	}
}

func (b *Backend) Update(status mender.DeploymentStatus, progress int) {
	st := UpdateStateStatus{
		State:    status,
		Progress: progress,
	}

	b.eventEmitter.OnUpdateStatus(st)

	if status == mender.Success || status == mender.Failure || status == mender.AlreadyInstalled {
		b.finishUpdate(status == mender.Success)
	}
}

func (b *Backend) Commit(success bool) {
	if b.update.Updating {
		b.waitCommit <- success
	}
}

func (b *Backend) MoveToNextState(move bool) {
	if b.update.Updating {
		b.waitCommit <- move
	}
}

func (b *Backend) NextState(status mender.DeploymentStatus) bool {
	b.eventEmitter.UpdateNextState(UpdateNextState{State: status})
	return <-b.waitCommit
}

func (b *Backend) Error(err error) {
	logger.Error("UpdateError", logging.String("error", err.Error()))
	b.eventEmitter.OnError(ErrUpdate)
}

func (b *Backend) StopUpdate() {
	if err := b.updater.StopUpdate(); err != nil {
		b.eventEmitter.OnError(ErrStopUpdate)
	}

	b.finishUpdate(false)
}

func (b *Backend) finishUpdate(success bool) {
	if b.update.Updating {
		b.update.Success = success
		b.update.Updating = false
		b.eventEmitter.OnUpdate(b.update)
		close(b.waitCommit)
	}
}
