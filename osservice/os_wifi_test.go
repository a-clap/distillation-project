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
	"os"

	"osservice"
	"osservice/mocks"
	"osservice/pkg/wifi"

	"github.com/golang/mock/gomock"
)

func (o *OsServiceSuite) wifiClient(msg string) *osservice.WifiClient {
	client, err := osservice.NewWifiClient(SrvTestHost, SrvTestPort, ClientTimeout)
	o.Require().Nil(err, msg)
	o.Require().NotNil(client, msg)
	return client
}

func (o *OsServiceSuite) TestWifi_APs() {
	args := []struct {
		name string
		aps  []wifi.AP
		err  error
	}{
		{
			name: "good",
			aps: []wifi.AP{
				{
					ID:   1,
					SSID: "123",
				},
				{
					ID:   2,
					SSID: "455",
				},
			},
			err: nil,
		},
		{
			name: "nil aps",
			aps:  nil,
			err:  nil,
		},
		{
			name: "empty aps",
			aps:  []wifi.AP{},
			err:  nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)

		waitAps := make(chan struct{})
		mockWifi.EXPECT().APs().DoAndReturn(func() ([]wifi.AP, error) {
			close(waitAps)
			return arg.aps, arg.err
		})

		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			aps, err := client.APs()
			o.waitFor(waitAps, "waitAps")

			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.ElementsMatch(arg.aps, aps, err)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestWifi_Disconnect() {
	args := []struct {
		name string
		err  error
	}{
		{
			name: "good",
			err:  nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Disconnect().Return(arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			err := client.Disconnect()
			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestWifi_Connect() {
	args := []struct {
		name string
		net  wifi.Network
		err  error
	}{
		{
			name: "good",
			net: wifi.Network{
				AP: wifi.AP{
					ID:   1,
					SSID: "123",
				},
				Password: "456",
			},
			err: nil,
		},
		{
			name: "err",
			err:  os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Connect(arg.net).Return(arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			err := client.Connect(arg.net)
			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
			}

			ctrl.Finish()
		})
	}
}

func (o *OsServiceSuite) TestWifi_Connected() {
	args := []struct {
		name      string
		retStatus wifi.Status
		err       error
	}{
		{
			name: "good",
			retStatus: wifi.Status{
				Connected: true,
				SSID:      "my ap",
			},
			err: nil,
		},
		{
			name: "good #2",
			retStatus: wifi.Status{
				Connected: false,
				SSID:      "",
			},
			err: nil,
		},
		{
			name: "err",
			retStatus: wifi.Status{
				Connected: false,
				SSID:      "",
			},
			err: os.ErrClosed,
		},
	}
	req := o.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(o.T())

		mockWifi := mocks.NewMockWifi(ctrl)
		mockWifi.EXPECT().Connected().Return(arg.retStatus, arg.err)
		opts := []osservice.Option{osservice.WithWifi(mockWifi)}

		new(TestServer).With(opts, func() {
			client := o.wifiClient(arg.name)
			status, err := client.Connected()
			if arg.err != nil {
				req.NotNil(err, arg.name)
				req.ErrorContains(err, arg.err.Error(), arg.name)
			} else {
				req.Nil(err, arg.name)
				req.EqualValues(arg.retStatus, status, arg.name)
			}

			ctrl.Finish()
		})
	}
}
