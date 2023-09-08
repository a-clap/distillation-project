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

package wifi

import (
	"errors"

	"osservice/pkg/wifi"
)

var (
	ErrNotImplemented = errors.New("not implemented")
	handler           = new(wifiHandler)
)

type Client interface {
	APs() ([]wifi.AP, error)
	Connected() (wifi.Status, error)
	Disconnect() error
	Connect(n wifi.Network) error
}

type Listener interface {
	OnWifiConfigChange()
}

type wifiHandler struct {
	Client Client
}

func Init(client Client) {
	handler.Client = client
}

func IsImplemented() bool {
	return handler.Client != nil
}

func isNotImplemented() bool {
	return !IsImplemented()
}

func Connect(ap, psk string) error {
	if isNotImplemented() {
		return ErrNotImplemented
	}
	w := wifi.Network{
		AP: wifi.AP{
			ID:   0,
			SSID: ap,
		},
		Password: psk,
	}

	return handler.Client.Connect(w)
}

func Connected() (connected bool, ap string, err error) {
	if isNotImplemented() {
		err = ErrNotImplemented
		return
	}
	s, err := handler.Client.Connected()
	if err != nil {
		return
	}
	return s.Connected, s.SSID, nil
}

func AP() ([]string, error) {
	if isNotImplemented() {
		return nil, ErrNotImplemented
	}
	nets, err := handler.Client.APs()
	if err != nil {
		return nil, err
	}
	aps := make([]string, len(nets))
	for i, net := range nets {
		aps[i] = net.SSID
	}
	return aps, nil
}

func Disconnect() error {
	if isNotImplemented() {
		return ErrNotImplemented
	}
	return handler.Client.Disconnect()
}
