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
	"fmt"
	"io"
	"time"
)

var (
	ErrLackOfInterface   = errors.New("lack of wireless interface")
	ErrNoSuchInterface   = errors.New("no such hardware interface")
	ErrAlreadyConnected  = errors.New("already connected")
	ErrNotConnected      = errors.New("not connected")
	ErrConnectionTimeout = errors.New("connection timeout")
	ErrConnectionError   = errors.New("connection")
	ErrDisconnectError   = errors.New("disconnect")
)

type Client interface {
	io.Closer
	SetScanTimeout(duration time.Duration)

	Status() (Status, error)
	Scan() ([]AP, error)
	ConnectWithEvents(n Network, events ...EventID) (<-chan Event, error)
	Disconnect() error
}

// Wireless is an interface based on go-wireless
// mainly created to make code at least a bit testable
type Wireless interface {
	Interfaces(basePath ...string) []string
	Client(iface string) (Client, error)
}

type Wifi struct {
	Wireless
	Interface  string
	client     Client
	interfaces []string
	extEvents  <-chan Event
	timeout    time.Duration
}

func NewWithInterface(iface Wireless) (*Wifi, error) {
	w := &Wifi{
		Wireless:   iface,
		interfaces: nil,
		Interface:  "",
		extEvents:  nil,
		timeout:    4 * time.Second,
	}

	w.interfaces = w.Wireless.Interfaces()
	if len(w.interfaces) == 0 {
		return nil, ErrLackOfInterface
	}

	var err error
	for _, chosen := range w.interfaces {
		w.Interface = chosen
		if err = w.Choose(w.Interface); err == nil {
			return w, err
		}
	}

	return nil, err
}

func New() (*Wifi, error) {
	w := newWireless()
	return NewWithInterface(w)
}

func (w *Wifi) newClient() (err error) {
	if w.client != nil {
		_ = w.client.Close()
	}

	w.client, err = w.Client(w.Interface)
	if err == nil {
		// by experience... this is a good value
		w.client.SetScanTimeout(w.timeout)
	}
	return
}

func (w *Wifi) Interfaces() []string {
	return w.interfaces
}

func (w *Wifi) SetCommandTimeout(timeout time.Duration) {
	w.timeout = timeout
}

func (w *Wifi) Choose(iface string) error {
	for _, exists := range w.interfaces {
		if exists == iface {
			w.Interface = iface
			return w.newClient()
		}
	}
	return ErrNoSuchInterface
}

// APs tries to get SSID list from wireless interface. It can block for a while
func (w *Wifi) APs() ([]AP, error) {
	aps, err := w.client.Scan()
	if err != nil {
		return nil, err
	}

	return aps, nil
}

func (w *Wifi) Connected() (Status, error) {
	s, err := w.client.Status()
	if err != nil {
		return s, err
	}
	return s, nil
}

func (w *Wifi) Disconnect() error {
	if c, err := w.Connected(); !c.Connected || err != nil {
		if err != nil {
			return err
		}
		return ErrNotConnected
	}

	defer func(w *Wifi) {
		w.extEvents = nil
	}(w)

	// Clear channel, as we will expect disconnect
	for len(w.extEvents) > 0 {
		<-w.extEvents
	}

	// What can we do about disconnect error? I think not much
	if err := w.client.Disconnect(); err != nil {
		return err
	}

	// If channel is nil, then Wifi was connected before creation of this object
	// can't handle disconnect event
	if w.extEvents == nil {
		return nil
	}

	ev, err := w.eventOrTimeout()
	if err != nil {
		return err
	}
	fmt.Println(ev)

	if ev.ID != Disconnected {
		return fmt.Errorf("%w: %d %s", ErrDisconnectError, ev.ID, ev.Message)
	}

	return nil
}

// Connect connects to specified Network
func (w *Wifi) Connect(n Network) error {
	s, err := w.client.Status()
	if err != nil {
		return err
	}

	// Already connected
	if s.Connected {
		return fmt.Errorf("%w to: %s", ErrAlreadyConnected, s.SSID)
	}

	events := []EventID{Connected, Disconnected, NetworkNotFound, AuthReject, OtherError}

	w.extEvents, err = w.client.ConnectWithEvents(n, events...)
	if err != nil {
		return err
	}

	ev, err := w.eventOrTimeout()
	if err != nil {
		return err
	}

	if ev.ID != Connected {
		return fmt.Errorf("%w: %s", ErrConnectionError, ev.Message)
	}

	return nil
}

func (w *Wifi) eventOrTimeout() (Event, error) {
	select {
	case ev := <-w.extEvents:
		return ev, nil
	case <-time.After(w.timeout):
		return Event{}, ErrConnectionTimeout
	}
}
