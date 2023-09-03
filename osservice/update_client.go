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
	"io"
	"sync"
	"time"

	"osservice/osproto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	_              Update = (*UpdateClient)(nil)
	ErrStreamEmpty        = errors.New("no stream")
)

type UpdateClient struct {
	timeout   time.Duration
	conn      *grpc.ClientConn
	client    osproto.UpdateClient
	stream    osproto.Update_UpdateClient
	streamMtx sync.Mutex
	finished  chan struct{}
}

func NewUpdateClient(addr string, port int, timeout time.Duration) (*UpdateClient, error) {
	conn, err := grpc.Dial(fmt.Sprintf("%v:%v", addr, port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &UpdateClient{
		timeout: timeout,
		conn:    conn,
		client:  osproto.NewUpdateClient(conn),
	}, nil
}

func (u *UpdateClient) ctx() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(u.timeout))
}

func (u *UpdateClient) ContinueUpdate() (bool, string) {
	ctx, cancel := u.ctx()
	defer cancel()

	updating, err := u.client.ContinueUpdate(ctx, &emptypb.Empty{})
	if err != nil {
		return false, ""
	}
	return updating.GetDuringUpdate().GetValue(), updating.GetUpdateName().GetValue()
}

func (u *UpdateClient) PullReleases() (bool, error) {
	ctx, cancel := u.ctx()
	defer cancel()

	newRelease, err := u.client.PullReleases(ctx, &emptypb.Empty{})
	if err != nil {
		return false, err
	}

	return newRelease.GetValue(), nil
}

func (u *UpdateClient) AvailableReleases() ([]string, error) {
	ctx, cancel := u.ctx()
	defer cancel()

	availableReleases, err := u.client.AvailableReleases(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	releases := make([]string, len(availableReleases.GetReleases()))
	for i, release := range availableReleases.GetReleases() {
		releases[i] = release.GetValue()
	}

	return releases, nil
}

func (u *UpdateClient) Update(artifactName string, callbacks UpdateCallbacks) error {
	return u.handleUpdate(artifactName, callbacks)
}

func (u *UpdateClient) handleUpdate(artifactName string, callbacks UpdateCallbacks) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2*time.Hour))
	var err error

	// Try to connect
	u.stream, err = u.client.Update(ctx)
	if err != nil {
		cancel()
		return err
	}

	// Send artifactName - if it is wrong, fail fast
	req := &osproto.UpdateRequest{
		ArtifactName: wrapperspb.String(artifactName),
		Stop:         nil,
		Continue:     nil,
	}

	if err := u.send(req); err != nil {
		cancel()
		return err
	}

	u.finished = make(chan struct{})
	// So far so good, run in background
	go func() {
		defer cancel()

		var response *osproto.UpdateResponse
		for {
			response, err = u.stream.Recv()
			if err != nil {
				break
			}

			// Handle response
			if response.Error != nil {
				err = errors.New(response.Error.GetValue())
				break
			}

			if response.NextState != nil {
				next := fromProtoUpdateState(*response.NextState)
				confirmed := callbacks.NextState(next)

				if err = u.send(&osproto.UpdateRequest{
					ArtifactName: nil,
					Stop:         nil,
					Continue:     wrapperspb.Bool(confirmed),
				}); err != nil {
					break
				}
				// Will we continue update?
				if !confirmed {
					break
				}
				continue
			}
			st := fromProtoUpdateState(response.State)
			callbacks.Update(st, int(response.Progress))

			if response.Finished != nil && response.Finished.GetValue() {
				break
			}
		}
		close(u.finished)
		_ = u.StopUpdate()

		// Nothing to do
		if err == nil || errors.Is(err, io.EOF) {
			return
		}

		// Something other, notify user
		if st, ok := status.FromError(err); ok {
			err = errors.New(st.Message())
			callbacks.Error(err)
		} else {
			// Shouldn't happen according to docs
			callbacks.Error(err)
		}
	}()

	return err
}

func (u *UpdateClient) StopUpdate() error {
	if u.stream != nil {
		closeReq := &osproto.UpdateRequest{
			ArtifactName: nil,
			Stop:         wrapperspb.Bool(true),
			Continue:     nil,
		}

		if err := u.send(closeReq); err != nil {
			return err
		}

		<-u.finished
		return u.stream.CloseSend()
	}
	return nil
}

func (u *UpdateClient) send(req *osproto.UpdateRequest) error {
	u.streamMtx.Lock()
	defer u.streamMtx.Unlock()

	if u.stream == nil {
		return ErrStreamEmpty
	}

	return u.stream.Send(req)
}
