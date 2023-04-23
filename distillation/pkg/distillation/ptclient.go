/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"context"
	"time"
	
	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
	"github.com/a-clap/embedded/pkg/restclient"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PTClient struct {
	addr    string
	timeout time.Duration
}

func NewPTClient(addr string, timeout time.Duration) *PTClient {
	return &PTClient{addr: addr, timeout: timeout}
}

func (p *PTClient) GetSensors() ([]PTConfig, error) {
	return restclient.Get[[]PTConfig, *Error](p.addr+RoutesGetPT, p.timeout)
}

func (p *PTClient) Configure(setConfig PTConfig) (PTConfig, error) {
	return restclient.Put[PTConfig, *Error](p.addr+RoutesConfigurePT, p.timeout, setConfig)
}

func (p *PTClient) Temperatures() ([]PTTemperature, error) {
	return restclient.Get[[]PTTemperature, *Error](p.addr+RoutesGetPTTemperatures, p.timeout)
}

type PTRPCClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  distillationproto.PTClient
}

func NewPTRPCClient(addr string, timeout time.Duration) (*PTRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &PTRPCClient{timeout: timeout, conn: conn, client: distillationproto.NewPTClient(conn)}, nil
}

func (g *PTRPCClient) Get() ([]PTConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.PTGet(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	confs := make([]PTConfig, len(got.Configs))
	for i, elem := range got.Configs {
		confs[i] = rpcToPTConfig(elem)
	}
	return confs, nil
}

func (g *PTRPCClient) Configure(setConfig PTConfig) (PTConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	set := ptConfigToRPC(&setConfig)
	got, err := g.client.PTConfigure(ctx, set)
	if err != nil {
		return PTConfig{}, err
	}
	setConfig = rpcToPTConfig(got)
	return setConfig, nil
}

func (g *PTRPCClient) Temperatures() ([]PTTemperature, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.PTGetTemperatures(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return rpcToPTTemperature(got), nil
}

func (g *PTRPCClient) Close() {
	_ = g.conn.Close()
}
