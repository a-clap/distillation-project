package osservice

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"osservice/osproto"
)

type Os struct {
	host string
	port int

	time  Time
	store Store
	net   Net
	wifi  Wifi

	osproto.UnimplementedTimeServer
	osproto.UnimplementedStoreServer
	osproto.UnimplementedNetServer
	osproto.UnimplementedWifiServer
	srv *grpc.Server
}

const (
	invalidPort = -1
)

func New(options ...Option) (*Os, error) {
	// Start with sane defaults
	o := &Os{
		host: "localhost",
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
	listener, err := net.Listen("tcp", fmt.Sprintf("%v:%v", o.host, o.port))
	if err != nil {
		return err
	}
	o.srv = grpc.NewServer()
	osproto.RegisterTimeServer(o.srv, o)
	osproto.RegisterStoreServer(o.srv, o)
	osproto.RegisterNetServer(o.srv, o)
	osproto.RegisterWifiServer(o.srv, o)

	return o.srv.Serve(listener)
}

func (o *Os) Stop() {
	o.srv.Stop()
}
