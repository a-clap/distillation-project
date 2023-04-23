package distillation

import (
	"context"
	"net"
	
	"github.com/a-clap/distillation/pkg/distillation/distillationproto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type RPC struct {
	distillationproto.UnimplementedHeaterServer
	distillationproto.UnimplementedGPIOServer
	distillationproto.UnimplementedDSServer
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
	distillationproto.RegisterDSServer(s, r)
	distillationproto.RegisterHeaterServer(s, r)
	
	r.Distillation.Run()
	defer r.Distillation.Close()
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

func (r *RPC) DSGet(ctx context.Context, e *empty.Empty) (*distillationproto.DSConfigs, error) {
	logger.Debug("DSGet")
	g := r.Distillation.DSHandler.GetSensors()
	
	configs := make([]*distillationproto.DSConfig, len(g))
	for i, elem := range g {
		configs[i] = dsConfigToRPC(&elem)
	}
	return &distillationproto.DSConfigs{Configs: configs}, nil
}

func (r *RPC) DSConfigure(ctx context.Context, config *distillationproto.DSConfig) (*distillationproto.DSConfig, error) {
	logger.Debug("DSConfigure")
	cfg := rpcToDSConfig(config)
	newCfg, err := r.Distillation.DSHandler.ConfigureSensor(cfg)
	if err != nil {
		return nil, err
	}
	return dsConfigToRPC(&newCfg), nil
}

func (r *RPC) DSGetTemperatures(ctx context.Context, e *empty.Empty) (*distillationproto.DSTemperatures, error) {
	logger.Debug("DSGetTemperatures")
	t := r.Distillation.DSHandler.Temperatures()
	return dsTemperatureToRPC(t), nil
}

func (r *RPC) HeaterGet(ctx context.Context, e *empty.Empty) (*distillationproto.HeaterConfigs, error) {
	logger.Debug("HeaterGet")
	g := r.Distillation.HeatersHandler.ConfigsGlobal()
	
	configs := make([]*distillationproto.HeaterConfig, len(g))
	for i, elem := range g {
		configs[i] = heaterConfigToRPC(&elem)
	}
	return &distillationproto.HeaterConfigs{Configs: configs}, nil
}

func (r *RPC) HeaterConfigure(ctx context.Context, config *distillationproto.HeaterConfig) (*distillationproto.HeaterConfig, error) {
	logger.Debug("HeaterConfigure")
	cfg := rpcToHeaterConfig(config)
	newCfg, err := r.Distillation.HeatersHandler.ConfigureGlobal(cfg)
	if err != nil {
		return nil, err
	}
	
	return heaterConfigToRPC(&newCfg), nil
}
