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

package osservice

import (
	"context"
	"errors"
	"fmt"
	"mender"
	"time"

	"osservice/osproto"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"go.uber.org/atomic"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var ErrForceStop = errors.New("client stopped update")

type UpdateCallbacks interface {
	mender.Callbacks
}

type Update interface {
	ContinueUpdate() (bool, string)
	PullReleases() (bool, error)
	AvailableReleases() ([]string, error)
	Update(artifactName string, callbacks UpdateCallbacks) error
	StopUpdate() error
}

var _ UpdateCallbacks = (*osUpdateCallbacks)(nil)

func (o *Os) PullReleases(context.Context, *empty.Empty) (*wrappers.BoolValue, error) {
	release, err := o.update.PullReleases()
	if err != nil {
		return nil, err
	}

	return wrapperspb.Bool(release), nil
}

func (o *Os) AvailableReleases(context.Context, *empty.Empty) (*osproto.Releases, error) {
	releases, err := o.update.AvailableReleases()
	if err != nil {
		return nil, err
	}

	protoReleases := &osproto.Releases{Releases: make([]*wrappers.StringValue, len(releases))}
	for i, release := range releases {
		protoReleases.Releases[i] = wrapperspb.String(release)
	}

	return protoReleases, nil
}

func (o *Os) ContinueUpdate(context.Context, *empty.Empty) (*osproto.UpdateInformation, error) {
	release, name := o.update.ContinueUpdate()

	return &osproto.UpdateInformation{
		DuringUpdate: wrapperspb.Bool(release),
		UpdateName:   wrapperspb.String(name),
	}, nil
}

func (o *Os) Update(server osproto.Update_UpdateServer) error {
	type updateStatus struct {
		status   mender.DeploymentStatus
		progress int
	}
	callbacks := &osUpdateCallbacks{}

	errs := make(chan error)
	update := make(chan updateStatus)
	nextState := make(chan mender.DeploymentStatus)
	nextStateConfirm := make(chan bool)

	callbacks.WithError = func(err error) {
		errs <- err
	}

	callbacks.WithUpdate = func(status mender.DeploymentStatus, progress int) {
		update <- updateStatus{
			status:   status,
			progress: progress,
		}
	}

	callbacks.WithNextState = func(status mender.DeploymentStatus) bool {
		nextState <- status
		return <-nextStateConfirm
	}

	// We need to get proper artifact name from client
	req, err := server.Recv()
	if err != nil {
		return err
	}

	if req.ArtifactName == nil {
		return fmt.Errorf("expected to receive artifact name")
	}

	err = o.update.Update(req.ArtifactName.GetValue(), callbacks)
	if err != nil {
		return err
	}

	fromClientFinished := make(chan struct{})
	waitingForResponse := atomic.NewBool(false)
	running := atomic.NewBool(true)

	// Handle message from client
	fromClient := func() {
		for {
			req, err := server.Recv()
			if !running.Load() {
				fmt.Println("closing client")
				break
			}

			if err != nil {
				errs <- err
				break
			}

			if req.Stop != nil {
				fmt.Println("req stop")
				_ = o.update.StopUpdate()
				errs <- ErrForceStop
				break
			}

			if waitingForResponse.CompareAndSwap(true, false) && req.Continue != nil {
				nextStateConfirm <- req.Continue.GetValue()
			}
		}

		close(fromClientFinished)
	}
	go fromClient()

	for running.Load() {
		select {
		case <-time.After(time.Hour):
		case err = <-errs:
			u := &osproto.UpdateResponse{
				Error: wrapperspb.String(err.Error()),
			}

			// Send error to client
			// It should reply, so we can close stream
			_ = server.Send(u)
			running.Store(false)

		case st := <-nextState:
			next := toProtoUpdateState(st)
			u := &osproto.UpdateResponse{
				State:     0,
				Progress:  0,
				NextState: &next,
			}

			waitingForResponse.Store(true)
			if err = server.Send(u); err != nil {
				running.Store(false)
			}

		case status := <-update:
			u := &osproto.UpdateResponse{
				State:    toProtoUpdateState(status.status),
				Progress: int32(status.progress),
			}
			finished := slices.Index([]mender.DeploymentStatus{mender.Success, mender.Failure, mender.AlreadyInstalled}, status.status) != -1
			if finished {
				u.Finished = wrapperspb.Bool(finished)
				running.Store(false)
			}

			if err = server.Send(u); err != nil {
				running.Store(false)
			}

		}
	}

	<-fromClientFinished
	return err
}

type osUpdateCallbacks struct {
	WithUpdate    func(status mender.DeploymentStatus, progress int)
	WithNextState func(status mender.DeploymentStatus) bool
	WithError     func(err error)
}

func (o *osUpdateCallbacks) Update(status mender.DeploymentStatus, progress int) {
	o.WithUpdate(status, progress)
}

func (o *osUpdateCallbacks) NextState(status mender.DeploymentStatus) bool {
	return o.WithNextState(status)
}

func (o *osUpdateCallbacks) Error(err error) {
	o.WithError(err)
}
