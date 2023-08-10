package osservice_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"osservice"
)

type OsServiceSuite struct {
	suite.Suite
}

const (
	SrvTestHost   = "localhost"
	SrvTestPort   = 50000
	ClientTimeout = time.Second
)

func TestOsService(t *testing.T) {
	suite.Run(t, new(OsServiceSuite))
}

type TestServer struct {
	srv             *osservice.Os
	started, finish chan struct{}
}

func (ts *TestServer) With(opts []osservice.Option, f func()) {
	ts.runWithOpts(opts)
	f()
	ts.stop()
}

func (ts *TestServer) runWithOpts(opts []osservice.Option) {
	var (
		defaultOpts = []osservice.Option{
			osservice.WithHost(SrvTestHost),
			osservice.WithPort(SrvTestPort),
		}
	)
	if ts.srv != nil {
		panic("you are trying to run TestServer twice")
	}

	var err error
	defaultOpts = append(defaultOpts, opts...)
	ts.srv, err = osservice.New(defaultOpts...)
	if err != nil {
		panic(fmt.Sprintln("failed to start server with opts:", opts))
	}
	ts.started = make(chan struct{})
	ts.finish = make(chan struct{})

	go func() {
		close(ts.started)
		_ = ts.srv.Run()
		close(ts.finish)
	}()

	<-ts.started

}

func (ts *TestServer) stop() {
	if ts.srv == nil {
		return
	}
	ts.srv.Stop()
	<-ts.finish
}
