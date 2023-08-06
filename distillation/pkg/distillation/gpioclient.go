/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"context"
	"time"

	"distillation/pkg/distillation/distillationproto"
	"embedded/pkg/restclient"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GPIOClient struct {
	addr    string
	timeout time.Duration
}

func NewGPIOClient(addr string, timeout time.Duration) *GPIOClient {
	return &GPIOClient{addr: addr, timeout: timeout}
}

func (h *GPIOClient) Get() ([]GPIOConfig, error) {
	return restclient.Get[[]GPIOConfig, *Error](h.addr+RoutesGetGPIO, h.timeout)
}

func (h *GPIOClient) Configure(setConfig GPIOConfig) (GPIOConfig, error) {
	return restclient.Put[GPIOConfig, *Error](h.addr+RoutesConfigureGPIO, h.timeout, setConfig)
}

type GPIORPCClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  distillationproto.GPIOClient
}

func NewGPIORPCClient(addr string, timeout time.Duration) (*GPIORPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &GPIORPCClient{timeout: timeout, conn: conn, client: distillationproto.NewGPIOClient(conn)}, nil
}

func (g *GPIORPCClient) Get() ([]GPIOConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.GPIOGet(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	confs := make([]GPIOConfig, len(got.Configs))
	for i, elem := range got.Configs {
		confs[i] = rpcToGPIOConfig(elem)
	}
	return confs, nil
}

func (g *GPIORPCClient) Configure(setConfig GPIOConfig) (GPIOConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	set := gpioConfigToRPC(&setConfig)
	got, err := g.client.GPIOConfigure(ctx, set)
	if err != nil {
		return GPIOConfig{}, err
	}
	setConfig = rpcToGPIOConfig(got)
	return setConfig, nil
}

func (g *GPIORPCClient) Close() {
	_ = g.conn.Close()
}
