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

package osservice_test

import (
	"fmt"
	"mender"
	"os"

	"osservice"
	"osservice/mocks"

	"github.com/golang/mock/gomock"
)

func (o *OsServiceSuite) updateClient() *osservice.UpdateClient {
	client, err := osservice.NewUpdateClient(SrvTestHost, SrvTestPort, ClientTimeout)
	o.Require().Nil(err)
	o.Require().NotNil(client)
	return client
}

func (o *OsServiceSuite) TestUpdate_HappyPath() {
	req := o.Require()

	ctrl := gomock.NewController(o.T())
	mockUpdate := mocks.NewMockUpdate(ctrl)

	releaseName := "release"

	var serverCallbacks osservice.UpdateCallbacks
	waitForCallbacks := make(chan struct{})

	mockUpdate.EXPECT().Update(releaseName, gomock.Any()).Do(func(artifactName string, cbs osservice.UpdateCallbacks) {
		serverCallbacks = cbs
		close(waitForCallbacks)
	}).Return(nil)

	clientCallbacks := mocks.NewMockUpdateCallbacks(ctrl)
	opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}

	handleStatus := func(status mender.DeploymentStatus, progress int) {
		waitForStatus := make(chan struct{})
		clientCallbacks.EXPECT().Update(status, progress).Do(func(any, any) {
			waitForStatus <- struct{}{}
		}).Times(1)

		serverCallbacks.Update(status, progress)
		o.waitFor(waitForStatus, "waitForStatus")
	}

	handleNext := func(nextStatus mender.DeploymentStatus) {
		waitForNext := make(chan struct{})
		clientCallbacks.EXPECT().NextState(nextStatus).DoAndReturn(func(s mender.DeploymentStatus) bool {
			close(waitForNext)
			return true
		}).Times(1)

		serverCallbacks.NextState(nextStatus)
		o.waitFor(waitForNext, "waitForNext")
	}

	handleStop := func(client *osservice.UpdateClient) error {
		waitStopUpdate := make(chan struct{})
		mockUpdate.EXPECT().StopUpdate().DoAndReturn(func() error {
			close(waitStopUpdate)
			return nil
		})

		waitError := make(chan struct{})
		clientCallbacks.EXPECT().Error(osservice.ErrForceStop).Do(func(any) {
			close(waitError)
		})

		err := client.StopUpdate()
		o.waitFor(waitError, "waitError")
		o.waitFor(waitStopUpdate, "waitStopUpdate")
		return err
	}

	new(TestServer).With(opts, func() {
		updateClient := o.updateClient()

		err := updateClient.Update(releaseName, clientCallbacks)
		req.Nil(err)

		o.waitFor(waitForCallbacks, "waitForCallbacks")
		// Downloading
		handleStatus(mender.Downloading, 1)
		handleStatus(mender.Downloading, 51)
		handleStatus(mender.Downloading, 100)
		// Next: PauseBeforeInstalling
		handleNext(mender.PauseBeforeInstalling)

		// Installing
		handleStatus(mender.Installing, 1)
		handleStatus(mender.Installing, 51)
		handleStatus(mender.Installing, 100)
		// Next: PauseBeforeRebooting
		handleNext(mender.PauseBeforeRebooting)

		// Rebooting
		handleStatus(mender.Rebooting, 1)
		handleStatus(mender.Rebooting, 51)
		handleStatus(mender.Rebooting, 100)
		// Next: PauseBeforeCommitting
		handleNext(mender.PauseBeforeCommitting)

		// Close
		err = handleStop(updateClient)
		req.Nil(err)

		ctrl.Finish()
	})
}

func (o *OsServiceSuite) TestUpdate_PassStatus() {
	req := o.Require()

	ctrl := gomock.NewController(o.T())
	mockUpdate := mocks.NewMockUpdate(ctrl)

	releaseName := "release"

	var serverCallbacks osservice.UpdateCallbacks
	waitForCallbacks := make(chan struct{})

	mockUpdate.EXPECT().Update(releaseName, gomock.Any()).Do(func(artifactName string, cbs osservice.UpdateCallbacks) {
		serverCallbacks = cbs
		close(waitForCallbacks)
	}).Return(nil)

	clientCallbacks := mocks.NewMockUpdateCallbacks(ctrl)
	opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}

	handleStatus := func(status mender.DeploymentStatus, maxProgress int) {
		for i := 0; i < maxProgress; i++ {
			waitForStatus := make(chan struct{})
			clientCallbacks.EXPECT().Update(status, i).Do(func(any, any) {
				waitForStatus <- struct{}{}
			})

			serverCallbacks.Update(status, i)
			o.waitFor(waitForStatus, "waitForStatus")
		}
	}

	handleStop := func(client *osservice.UpdateClient) error {
		waitStopUpdate := make(chan struct{})
		mockUpdate.EXPECT().StopUpdate().DoAndReturn(func() error {
			close(waitStopUpdate)
			return nil
		})

		waitError := make(chan struct{})
		clientCallbacks.EXPECT().Error(osservice.ErrForceStop).Do(func(any) {
			close(waitError)
		})

		err := client.StopUpdate()
		o.waitFor(waitError, "waitError")
		o.waitFor(waitStopUpdate, "waitStopUpdate")
		return err
	}

	new(TestServer).With(opts, func() {
		updateClient := o.updateClient()

		err := updateClient.Update(releaseName, clientCallbacks)
		req.Nil(err)

		o.waitFor(waitForCallbacks, "waitForCallbacks")

		// Pass different statuses
		times := 10
		handleStatus(mender.Downloading, times)
		handleStatus(mender.PauseBeforeInstalling, times)
		handleStatus(mender.Installing, times)
		handleStatus(mender.PauseBeforeRebooting, times)
		handleStatus(mender.Rebooting, times)
		handleStatus(mender.PauseBeforeCommitting, times)
		handleStatus(mender.Success, times)
		handleStatus(mender.Failure, times)
		handleStatus(mender.AlreadyInstalled, times)

		err = handleStop(updateClient)
		req.Nil(err)
		ctrl.Finish()
	})
}

func (o *OsServiceSuite) TestUpdate_StopUpdate() {
	req := o.Require()
	releaseName := "release"

	ctrl := gomock.NewController(o.T())

	waitUpdate := make(chan struct{})
	mockUpdate := mocks.NewMockUpdate(ctrl)
	mockUpdate.EXPECT().Update(releaseName, gomock.Any()).Do(func(any, any) {
		close(waitUpdate)
	}).Return(nil)

	clientCallbacks := mocks.NewMockUpdateCallbacks(ctrl)
	opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}

	handleStop := func(client *osservice.UpdateClient) error {
		waitStopUpdate := make(chan struct{})
		mockUpdate.EXPECT().StopUpdate().DoAndReturn(func() error {
			close(waitStopUpdate)
			return nil
		})

		waitError := make(chan struct{})
		clientCallbacks.EXPECT().Error(osservice.ErrForceStop).Do(func(any) {
			close(waitError)
		})

		err := client.StopUpdate()
		o.waitFor(waitError, "waitError")
		o.waitFor(waitStopUpdate, "waitStopUpdate")
		return err
	}

	new(TestServer).With(opts, func() {
		updateClient := o.updateClient()

		err := updateClient.Update(releaseName, clientCallbacks)
		req.Nil(err)

		o.waitFor(waitUpdate, "waitUpdate")
		err = handleStop(updateClient)

		req.Nil(err)
		ctrl.Finish()
	})
}

func (o *OsServiceSuite) TestUpdate_UpdateHandleError() {
	ctrl := gomock.NewController(o.T())
	mockUpdate := mocks.NewMockUpdate(ctrl)

	releaseName := "release"

	var serverCallbacks osservice.UpdateCallbacks
	waitForCallbacks := make(chan struct{})

	mockUpdate.EXPECT().Update(releaseName, gomock.Any()).Do(func(artifactName string, cbs osservice.UpdateCallbacks) {
		serverCallbacks = cbs
		close(waitForCallbacks)
	}).Return(nil)

	updateError := fmt.Errorf("updateError")
	clientCallbacks := mocks.NewMockUpdateCallbacks(ctrl)

	waitForError := make(chan struct{})
	clientCallbacks.EXPECT().Error(updateError).Do(func(error) {
		close(waitForError)
	})

	opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}
	req := o.Require()
	new(TestServer).With(opts, func() {
		updateClient := o.updateClient()

		err := updateClient.Update(releaseName, clientCallbacks)
		req.Nil(err)

		o.waitFor(waitForCallbacks, "waitForCallbacks")

		serverCallbacks.Error(updateError)

		o.waitFor(waitForError, "waitForError")

		ctrl.Finish()
	})
}

func (o *OsServiceSuite) TestUpdate_AvailableReleases() {
	args := []struct {
		name              string
		availableReleases []string
		err               error
	}{
		{
			name:              "all good",
			availableReleases: []string{"1"},
			err:               nil,
		},
		{
			name:              "all good #2",
			availableReleases: []string{"1", "2", "3"},
			err:               nil,
		},
		{
			name:              "all good #3",
			availableReleases: []string{},
			err:               nil,
		},
		{
			name:              "error",
			availableReleases: []string{"nope"},
			err:               os.ErrExist,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockUpdate := mocks.NewMockUpdate(ctrl)
		waitForAvailableReleases := make(chan struct{})
		mockUpdate.EXPECT().AvailableReleases().DoAndReturn(func() ([]string, error) {
			close(waitForAvailableReleases)
			return arg.availableReleases, arg.err
		})

		opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}

		new(TestServer).With(opts, func() {
			updateClient := o.updateClient()

			newRelease, err := updateClient.AvailableReleases()
			o.waitFor(waitForAvailableReleases, "waitForAvailableReleases")
			if arg.err != nil {
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.ElementsMatch(newRelease, arg.availableReleases, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestUpdate_PullReleases() {
	args := []struct {
		name       string
		newRelease bool
		err        error
	}{
		{
			name:       "all good",
			newRelease: true,
			err:        nil,
		},
		{
			name:       "all good",
			newRelease: false,
			err:        nil,
		},
		{
			name:       "error",
			newRelease: false,
			err:        os.ErrExist,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockUpdate := mocks.NewMockUpdate(ctrl)
		waitForPull := make(chan struct{})
		mockUpdate.EXPECT().PullReleases().DoAndReturn(func() (bool, error) {
			close(waitForPull)
			return arg.newRelease, arg.err
		})

		opts := []osservice.Option{osservice.WithUpdate(mockUpdate)}

		new(TestServer).With(opts, func() {
			updateClient := o.updateClient()

			newRelease, err := updateClient.PullReleases()
			o.waitFor(waitForPull, "waitForPull")

			if arg.err != nil {
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.EqualValues(newRelease, arg.newRelease, arg.name)
			}

			ctrl.Finish()
		})
	}
}
