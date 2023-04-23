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

type DSClient struct {
	addr    string
	timeout time.Duration
}

func NewDSClient(addr string, timeout time.Duration) *DSClient {
	return &DSClient{addr: addr, timeout: timeout}
}

func (d *DSClient) GetSensors() ([]DSConfig, error) {
	return restclient.Get[[]DSConfig, *Error](d.addr+RoutesGetDS, d.timeout)
}

func (d *DSClient) Configure(setConfig DSConfig) (DSConfig, error) {
	return restclient.Put[DSConfig, *Error](d.addr+RoutesGetDS, d.timeout, setConfig)
}

func (d *DSClient) Temperatures() ([]DSTemperature, error) {
	return restclient.Get[[]DSTemperature, *Error](d.addr+RoutesGetDSTemperatures, d.timeout)
}

type DSRPCClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  distillationproto.DSClient
}

func NewDSRPCClient(addr string, timeout time.Duration) (*DSRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &DSRPCClient{timeout: timeout, conn: conn, client: distillationproto.NewDSClient(conn)}, nil
}

func (g *DSRPCClient) Get() ([]DSConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.DSGet(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	confs := make([]DSConfig, len(got.Configs))
	for i, elem := range got.Configs {
		confs[i] = rpcToDSConfig(elem)
	}
	return confs, nil
}

func (g *DSRPCClient) Configure(setConfig DSConfig) (DSConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	set := dsConfigToRPC(&setConfig)
	got, err := g.client.DSConfigure(ctx, set)
	if err != nil {
		return DSConfig{}, err
	}
	setConfig = rpcToDSConfig(got)
	return setConfig, nil
}

func (g *DSRPCClient) Temperatures() ([]DSTemperature, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.DSGetTemperatures(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	return rpcToDSTemperature(got), nil
}

func (g *DSRPCClient) Close() {
	_ = g.conn.Close()
}
