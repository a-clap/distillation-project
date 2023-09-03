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
	"fmt"
	"net"

	"google.golang.org/grpc"

	"osservice/osproto"
)

type Os struct {
	port int

	time   Time
	store  Store
	net    Net
	wifi   Wifi
	update Update

	osproto.UnimplementedTimeServer
	osproto.UnimplementedStoreServer
	osproto.UnimplementedNetServer
	osproto.UnimplementedWifiServer
	osproto.UnimplementedUpdateServer
	srv *grpc.Server
}

const (
	invalidPort = -1
)

func New(options ...Option) (*Os, error) {
	// Start with sane defaults
	o := &Os{
		port: invalidPort,
		time: timeOs{},
		net:  netOs{},
		wifi: newWifiOs(),
	}

	for _, option := range options {
		option(o)
	}

	if err := o.verify(); err != nil {
		return nil, err
	}

	return o, nil
}

func (o *Os) verify() error {
	if o.port == invalidPort {
		return fmt.Errorf(`you must specify port with "WithPort"`)
	}

	return nil
}

func (o *Os) Run() error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", o.port))
	if err != nil {
		return err
	}
	o.srv = grpc.NewServer()
	osproto.RegisterTimeServer(o.srv, o)
	osproto.RegisterStoreServer(o.srv, o)
	osproto.RegisterNetServer(o.srv, o)
	osproto.RegisterWifiServer(o.srv, o)
	osproto.RegisterUpdateServer(o.srv, o)

	return o.srv.Serve(listener)
}

func (o *Os) Stop() {
	o.srv.GracefulStop()
}
