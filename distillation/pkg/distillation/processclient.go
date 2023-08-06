/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"context"
	"strconv"
	"strings"
	"time"

	"distillation/pkg/distillation/distillationproto"
	"distillation/pkg/process"
	"embedded/pkg/restclient"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProcessClient struct {
	addr    string
	timeout time.Duration
}

func NewProcessClient(addr string, timeout time.Duration) *ProcessClient {
	return &ProcessClient{addr: addr, timeout: timeout}
}

func (h *ProcessClient) GetPhaseCount() (ProcessPhaseCount, error) {
	return restclient.Get[ProcessPhaseCount, *Error](h.addr+RoutesProcessPhases, h.timeout)
}

func (h *ProcessClient) GetPhaseConfig(phaseNumber int) (ProcessPhaseConfig, error) {
	p := strconv.FormatInt(int64(phaseNumber), 10)
	addr := strings.Replace(RoutesProcessConfigPhase, ":id", p, 1)
	return restclient.Get[ProcessPhaseConfig, *Error](h.addr+addr, h.timeout)
}

func (h *ProcessClient) ConfigurePhaseCount(count ProcessPhaseCount) (ProcessPhaseCount, error) {
	return restclient.Put[ProcessPhaseCount, *Error](h.addr+RoutesProcessPhases, h.timeout, count)
}

func (h *ProcessClient) ConfigurePhase(phaseNumber int, setConfig ProcessPhaseConfig) (ProcessPhaseConfig, error) {
	p := strconv.FormatInt(int64(phaseNumber), 10)
	addr := strings.Replace(RoutesProcessConfigPhase, ":id", p, 1)
	return restclient.Put[ProcessPhaseConfig, *Error](h.addr+addr, h.timeout, setConfig)
}

func (h *ProcessClient) ValidateConfig() (ProcessConfigValidation, error) {
	return restclient.Get[ProcessConfigValidation, *Error](h.addr+RoutesProcessConfigValidate, h.timeout)
}

func (h *ProcessClient) ConfigureProcess(cfg ProcessConfig) (ProcessConfig, error) {
	return restclient.Put[ProcessConfig, *Error](h.addr+RoutesProcess, h.timeout, cfg)
}

func (h *ProcessClient) Status() (ProcessStatus, error) {
	return restclient.Get[ProcessStatus, *Error](h.addr+RoutesProcessStatus, h.timeout)
}

type ProcessRPCClient struct {
	timeout time.Duration
	conn    *grpc.ClientConn
	client  distillationproto.ProcessClient
}

func NewProcessRPCClient(addr string, timeout time.Duration) (*ProcessRPCClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &ProcessRPCClient{timeout: timeout, conn: conn, client: distillationproto.NewProcessClient(conn)}, nil
}

func (p *ProcessRPCClient) GetPhaseCount() (ProcessPhaseCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	cnt, err := p.client.GetPhaseCount(ctx, &emptypb.Empty{})
	if err != nil {
		return ProcessPhaseCount{}, err
	}
	return rpcToProcessPhaseCount(cnt), err
}

func (p *ProcessRPCClient) GetPhaseConfig(phaseNumber int) (ProcessPhaseConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	cfg, err := p.client.GetPhaseConfig(ctx, &distillationproto.PhaseNumber{Number: int32(phaseNumber)})
	if err != nil {
		return ProcessPhaseConfig{}, err
	}
	return rpcToProcessPhaseConfig(cfg), nil

}

func (p *ProcessRPCClient) ConfigurePhaseCount(count ProcessPhaseCount) (ProcessPhaseCount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	cfg, err := p.client.ConfigurePhaseCount(ctx, &distillationproto.ProcessPhaseCount{Count: int32(count.PhaseNumber)})
	if err != nil {
		return ProcessPhaseCount{}, err
	}
	return ProcessPhaseCount{PhaseNumber: uint(cfg.Count)}, nil
}

func (p *ProcessRPCClient) ConfigurePhase(phaseNumber int, setConfig ProcessPhaseConfig) (ProcessPhaseConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()
	c := processPhaseConfigToRpc(phaseNumber, setConfig)

	cfg, err := p.client.ConfigurePhase(ctx, c)
	if err != nil {
		return ProcessPhaseConfig{}, err
	}
	return rpcToProcessPhaseConfig(cfg), nil
}

func (p *ProcessRPCClient) ConfigureGlobalGPIO(gpios []process.GPIOConfig) ([]process.GPIOConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	rpcConf := &distillationproto.GlobalGPIOConfig{
		Configs: make([]*distillationproto.GPIOPhaseConfig, len(gpios)),
	}
	for i, gp := range gpios {
		rpcConf.Configs[i] = gpioPhaseConfigToRpc(gp)
	}

	cfg, err := p.client.ConfigureGlobalGPIO(ctx, rpcConf)
	if err != nil {
		return nil, err
	}

	gpios = make([]process.GPIOConfig, len(cfg.Configs))
	for i, gp := range cfg.Configs {
		gpios[i] = rpcToGPIOPhaseConfig(gp)
	}

	return gpios, nil
}

func (p *ProcessRPCClient) ValidateConfig() (ProcessConfigValidation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	v, err := p.client.ValidateConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return ProcessConfigValidation{}, err
	}
	return ProcessConfigValidation{Valid: v.Valid, Error: v.Error}, nil
}

func (p *ProcessRPCClient) ConfigureProcess(cfg ProcessConfig) (ProcessConfig, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	conf := processConfigToRpc(cfg)
	conf, err := p.client.EnableProcess(ctx, conf)
	if err != nil {
		return ProcessConfig{}, err
	}
	return rpcToProcessConfig(conf), nil
}

func (p *ProcessRPCClient) Status() (ProcessStatus, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	status, err := p.client.Status(ctx, &emptypb.Empty{})
	if err != nil {
		return ProcessStatus{}, err
	}
	return rpcToProcessStatus(status), nil
}

func (p *ProcessRPCClient) GlobalConfig() (process.Config, error) {
	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	defer cancel()

	global, err := p.client.GetGlobalConfig(ctx, &emptypb.Empty{})
	if err != nil {
		return process.Config{}, err
	}

	cfg := process.Config{
		PhaseNumber: uint(global.Count),
		Phases:      make([]process.PhaseConfig, len(global.PhaseConfig)),
		GlobalGPIO:  make([]process.GPIOConfig, len(global.GlobalGPIOConfig)),
		Sensors:     global.Sensors,
	}

	for i, gpio := range global.GlobalGPIOConfig {
		cfg.GlobalGPIO[i] = process.GPIOConfig{
			Enabled:    gpio.Enabled,
			ID:         gpio.ID,
			SensorID:   gpio.SensorID,
			TLow:       float64(gpio.TLow),
			THigh:      float64(gpio.THigh),
			Hysteresis: float64(gpio.Hysteresis),
			Inverted:   gpio.Inverted,
		}
	}

	for i, elem := range global.PhaseConfig {
		cfg.Phases[i] = rpcToProcessPhaseConfig(elem).PhaseConfig
	}

	return cfg, nil
}
