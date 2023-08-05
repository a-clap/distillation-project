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

type HeaterClient struct {
	addr    string
	timeout time.Duration
}

func NewHeaterClient(addr string, timeout time.Duration) *HeaterClient {
	return &HeaterClient{addr: addr, timeout: timeout}
}

func (h *HeaterClient) GetEnabled() ([]HeaterConfig, error) {
	return restclient.Get[[]HeaterConfig, *Error](h.addr+RoutesGetEnabledHeaters, h.timeout)
}

func (h *HeaterClient) GetAll() ([]HeaterConfigGlobal, error) {
	return restclient.Get[[]HeaterConfigGlobal, *Error](h.addr+RoutesGetAllHeaters, h.timeout)
}

func (h *HeaterClient) Enable(setConfig HeaterConfigGlobal) (HeaterConfigGlobal, error) {
	return restclient.Put[HeaterConfigGlobal, *Error](h.addr+RoutesEnableHeater, h.timeout, setConfig)
}

func (h *HeaterClient) Configure(setConfig HeaterConfig) (HeaterConfig, error) {
	return restclient.Put[HeaterConfig, *Error](h.addr+RoutesConfigureHeater, h.timeout, setConfig)
}

type HeaterRPCClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  distillationproto.HeaterClient
}

func NewHeaterRPCCLient(addr string, timeout time.Duration) (*HeaterRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &HeaterRPCClient{timeout: timeout, conn: conn, client: distillationproto.NewHeaterClient(conn)}, nil
}

func (g *HeaterRPCClient) Get() ([]HeaterConfigGlobal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	got, err := g.client.HeaterGet(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}
	confs := make([]HeaterConfigGlobal, len(got.Configs))
	for i, elem := range got.Configs {
		confs[i] = rpcToHeaterConfig(elem)
	}
	return confs, nil
}

func (g *HeaterRPCClient) Configure(setConfig HeaterConfigGlobal) (HeaterConfigGlobal, error) {
	ctx, cancel := context.WithTimeout(context.Background(), g.timeout)
	defer cancel()
	set := heaterConfigToRPC(&setConfig)
	got, err := g.client.HeaterConfigure(ctx, set)
	if err != nil {
		return HeaterConfigGlobal{}, err
	}
	setConfig = rpcToHeaterConfig(got)
	return setConfig, nil
}

func (g *HeaterRPCClient) Close() {
	_ = g.conn.Close()
}
