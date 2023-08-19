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
	"os"
	"time"

	"osservice"
	"osservice/mocks"

	"github.com/golang/mock/gomock"
)

func (o *OsServiceSuite) timeClient(msg string) *osservice.TimeClient {
	client, err := osservice.NewTimeClient(SrvTestHost, SrvTestPort, ClientTimeout)
	o.Require().Nil(err, msg)
	o.Require().NotNil(client, msg)
	return client
}

func (o *OsServiceSuite) TestTime_SetNTP() {
	args := []struct {
		name string
		set  bool
		err  error
	}{
		{
			name: "all good",
			set:  true,
			err:  nil,
		},
		{
			name: "all good",
			set:  false,
			err:  nil,
		},
		{
			name: "error",
			set:  false,
			err:  os.ErrExist,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockTime := mocks.NewMockTime(ctrl)
		waitSet := make(chan struct{})
		mockTime.EXPECT().SetNTP(arg.set).DoAndReturn(func(any) error {
			close(waitSet)
			return arg.err
		})
		opts := []osservice.Option{osservice.WithTime(mockTime)}

		new(TestServer).With(opts, func() {
			timeClient := o.timeClient(arg.name)

			err := timeClient.SetNTP(arg.set)
			o.waitFor(waitSet, "waitSet")

			if arg.err != nil {
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
			}
			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestTime_SetNow() {
	args := []struct {
		name string
		set  time.Time
		err  error
	}{
		{
			name: "all good",
			set:  time.UnixMilli(1),
			err:  nil,
		},
		{
			name: "error",
			set:  time.UnixMilli(1),
			err:  os.ErrExist,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockTime := mocks.NewMockTime(ctrl)
		waitSetNow := make(chan struct{})
		mockTime.EXPECT().SetNow(arg.set.UTC()).DoAndReturn(func(any) error {
			close(waitSetNow)
			return arg.err
		})

		opts := []osservice.Option{osservice.WithTime(mockTime)}

		new(TestServer).With(opts, func() {
			timeClient := o.timeClient(arg.name)

			err := timeClient.SetNow(arg.set.UTC())
			o.waitFor(waitSetNow, "waitSetNow")

			if arg.err != nil {
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestTime_NTP() {
	args := []struct {
		name    string
		timeout time.Duration
		in      struct {
			ntp bool
			err error
		}
		expected struct {
			ntp bool
			err error
		}
	}{
		{
			name:    "all good - disabled",
			timeout: time.Second,
			in: struct {
				ntp bool
				err error
			}{ntp: false, err: nil},
			expected: struct {
				ntp bool
				err error
			}{ntp: false, err: nil},
		},
		{
			name:    "all good - enabled",
			timeout: time.Second,
			in: struct {
				ntp bool
				err error
			}{ntp: true, err: nil},
			expected: struct {
				ntp bool
				err error
			}{ntp: true, err: nil},
		},
		{
			name:    "err",
			timeout: time.Second,
			in: struct {
				ntp bool
				err error
			}{ntp: true, err: os.ErrClosed},
			expected: struct {
				ntp bool
				err error
			}{ntp: false, err: os.ErrClosed},
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockTime := mocks.NewMockTime(ctrl)
		waitNTP := make(chan struct{})
		mockTime.EXPECT().NTP().DoAndReturn(func() (bool, error) {
			close(waitNTP)
			return arg.in.ntp, arg.in.err
		})

		opts := []osservice.Option{osservice.WithTime(mockTime)}

		new(TestServer).With(opts, func() {
			timeClient := o.timeClient(arg.name)

			ntp, err := timeClient.NTP()
			o.waitFor(waitNTP, "waitNTP")

			if arg.expected.err != nil {
				req.Equal(arg.expected.ntp, ntp, arg.name)
				req.ErrorContains(err, arg.expected.err.Error(), arg.name)
			} else {
				req.Equal(arg.expected.ntp, ntp, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestTime_Now() {
	args := []struct {
		name    string
		timeout time.Duration
		in      struct {
			now time.Time
			err error
		}
		expected struct {
			now time.Time
			err error
		}
	}{
		{
			name:    "all good",
			timeout: time.Second,
			in: struct {
				now time.Time
				err error
			}{now: time.UnixMilli(1), err: nil},
			expected: struct {
				now time.Time
				err error
			}{now: time.UnixMilli(1), err: nil},
		},
		{
			name:    "all good #2",
			timeout: time.Second,
			in: struct {
				now time.Time
				err error
			}{now: time.UnixMilli(11212), err: nil},
			expected: struct {
				now time.Time
				err error
			}{now: time.UnixMilli(11212), err: nil},
		},
		{
			name:    "all good #3",
			timeout: time.Second,
			in: struct {
				now time.Time
				err error
			}{now: time.Date(2023, 9, 23, 13, 35, 17, 23, time.Local), err: nil},
			expected: struct {
				now time.Time
				err error
			}{now: time.Date(2023, 9, 23, 13, 35, 17, 23, time.Local), err: nil},
		},
		{
			name:    "return error",
			timeout: time.Second,
			in: struct {
				now time.Time
				err error
			}{now: time.Date(2023, 9, 23, 13, 35, 17, 23, time.Local), err: os.ErrClosed},
			expected: struct {
				now time.Time
				err error
			}{now: time.Time{}, err: os.ErrClosed},
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockTime := mocks.NewMockTime(ctrl)
		mockTime.EXPECT().Now().Return(arg.in.now, arg.in.err)
		opts := []osservice.Option{osservice.WithTime(mockTime)}

		new(TestServer).With(opts, func() {
			timeClient := o.timeClient(arg.name)

			now, err := timeClient.Now()
			if arg.expected.err != nil {
				req.Equal(arg.expected.now, now, arg.name)
				req.ErrorContains(err, arg.expected.err.Error(), arg.name)
			} else {
				req.Nil(err, fmt.Sprint(err, ": ", arg.name))
				req.True(arg.expected.now.Compare(now) == 0, arg.name)
			}

			ctrl.Finish()
		})
	}
}
