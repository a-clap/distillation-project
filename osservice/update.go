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
	"mender"
)

var _ Update = (*updateOs)(nil)

type updateOs struct {
	client    *mender.Client
	connected bool
}

func newUpdateOs(c *mender.Client) *updateOs {
	return &updateOs{
		client: c,
	}
}

func (u *updateOs) lazyInit() error {
	if !u.connected {
		if err := u.client.Connect(); err != nil {
			return err
		}

		if err := u.client.UpdateInventory(); err != nil {
			return err
		}
		u.connected = true
	}

	return nil
}

func (u *updateOs) ContinueUpdate() (bool, string) {
	if err := u.lazyInit(); err != nil {
		return false, ""
	}
	return u.client.ContinueUpdate()
}

func (u *updateOs) PullReleases() (bool, error) {
	if err := u.lazyInit(); err != nil {
		return false, err
	}
	return u.client.PullReleases()
}

func (u *updateOs) AvailableReleases() ([]string, error) {
	if err := u.lazyInit(); err != nil {
		return nil, err
	}
	return u.client.AvailableReleases(), nil
}

func (u *updateOs) Update(artifactName string, callbacks UpdateCallbacks) error {
	if err := u.lazyInit(); err != nil {
		return err
	}
	u.client.Callbacks = callbacks
	return u.client.Update(artifactName)
}

func (u *updateOs) StopUpdate() error {
	if err := u.lazyInit(); err != nil {
		return err
	}
	return u.client.StopUpdate()
}
