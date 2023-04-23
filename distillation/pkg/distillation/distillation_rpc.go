package distillation

import (
	"context"
	"net"
	
	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type RPC struct {
	distillationproto.UnimplementedGPIOServer
	*Distillation
}

func NewRPC(options ...Option) (*RPC, error) {
	r := &RPC{}
	
	d, err := New(options...)
	if err != nil {
		return nil, err
	}
	r.Distillation = d
	
	return r, nil
}

func (r *RPC) Run() error {
	listener, err := net.Listen("tcp", r.Distillation.url)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	distillationproto.RegisterGPIOServer(s, r)
	
	return s.Serve(listener)
}

func (r *RPC) Close() {
}

func (r *RPC) GPIOGet(ctx context.Context, empty *empty.Empty) (*distillationproto.GPIOConfigs, error) {
	logger.Debug("GPIOGet")
	g := r.Distillation.GPIOHandler.Config()
	configs := make([]*distillationproto.GPIOConfig, len(g))
	for i, elem := range g {
		configs[i] = gpioConfigToRPC(&elem)
	}
	return &distillationproto.GPIOConfigs{Configs: configs}, nil
}

func (r *RPC) GPIOConfigure(ctx context.Context, cfg *distillationproto.GPIOConfig) (*distillationproto.GPIOConfig, error) {
	logger.Debug("GPIOConfigure")
	config := rpcToGPIOConfig(cfg)
	newCfg, err := r.Distillation.GPIOHandler.Configure(config)
	if err != nil {
		return nil, err
	}
	cfg = gpioConfigToRPC(&newCfg)
	return cfg, nil
}
