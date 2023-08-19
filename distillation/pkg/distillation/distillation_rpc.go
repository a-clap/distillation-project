package distillation

import (
	"context"
	"errors"
	"net"

	"distillation/pkg/distillation/distillationproto"
	"distillation/pkg/process"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type RPC struct {
	url string
	distillationproto.UnimplementedHeaterServer
	distillationproto.UnimplementedGPIOServer
	distillationproto.UnimplementedDSServer
	distillationproto.UnimplementedPTServer
	distillationproto.UnimplementedProcessServer
	*Distillation
}

func NewRPC(url string, options ...Option) (*RPC, error) {
	r := &RPC{url: url}

	d, err := New(options...)
	if err != nil {
		return nil, err
	}
	r.Distillation = d

	return r, nil
}

func (r *RPC) Run() error {
	listener, err := net.Listen("tcp", r.url)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	distillationproto.RegisterGPIOServer(s, r)
	distillationproto.RegisterDSServer(s, r)
	distillationproto.RegisterHeaterServer(s, r)
	distillationproto.RegisterPTServer(s, r)
	distillationproto.RegisterProcessServer(s, r)

	r.Distillation.Run()
	defer r.Distillation.Close()
	return s.Serve(listener)
}

func (r *RPC) Close() {
}

func (r *RPC) GPIOGet(ctx context.Context, empty *empty.Empty) (*distillationproto.GPIOConfigs, error) {
	g := r.Distillation.GPIOHandler.Config()
	configs := make([]*distillationproto.GPIOConfig, len(g))
	for i, elem := range g {
		configs[i] = gpioConfigToRPC(&elem)
	}
	return &distillationproto.GPIOConfigs{Configs: configs}, nil
}

func (r *RPC) GPIOConfigure(ctx context.Context, cfg *distillationproto.GPIOConfig) (*distillationproto.GPIOConfig, error) {
	config := rpcToGPIOConfig(cfg)
	newCfg, err := r.Distillation.GPIOHandler.Configure(config)
	if err != nil {
		return nil, err
	}
	cfg = gpioConfigToRPC(&newCfg)
	return cfg, nil
}

func (r *RPC) DSGet(ctx context.Context, e *empty.Empty) (*distillationproto.DSConfigs, error) {
	g := r.Distillation.DSHandler.GetSensors()

	configs := make([]*distillationproto.DSConfig, len(g))
	for i, elem := range g {
		configs[i] = dsConfigToRPC(&elem)
	}
	return &distillationproto.DSConfigs{Configs: configs}, nil
}

func (r *RPC) DSConfigure(ctx context.Context, config *distillationproto.DSConfig) (*distillationproto.DSConfig, error) {
	cfg := rpcToDSConfig(config)
	newCfg, err := r.Distillation.DSHandler.ConfigureSensor(cfg)
	if err != nil {
		return nil, err
	}
	return dsConfigToRPC(&newCfg), nil
}

func (r *RPC) DSGetTemperatures(ctx context.Context, e *empty.Empty) (*distillationproto.DSTemperatures, error) {
	t := r.Distillation.DSHandler.Temperatures()
	return dsTemperatureToRPC(t), nil
}

func (r *RPC) HeaterGet(ctx context.Context, e *empty.Empty) (*distillationproto.HeaterConfigs, error) {
	g := r.Distillation.HeatersHandler.ConfigsGlobal()

	configs := make([]*distillationproto.HeaterConfig, len(g))
	for i, elem := range g {
		configs[i] = heaterConfigToRPC(&elem)
	}
	return &distillationproto.HeaterConfigs{Configs: configs}, nil
}

func (r *RPC) HeaterConfigure(ctx context.Context, config *distillationproto.HeaterConfig) (*distillationproto.HeaterConfig, error) {
	cfg := rpcToHeaterConfig(config)
	newCfg, err := r.Distillation.HeatersHandler.ConfigureGlobal(cfg)
	if err != nil {
		return nil, err
	}

	return heaterConfigToRPC(&newCfg), nil
}

func (r *RPC) PTGet(ctx context.Context, e *empty.Empty) (*distillationproto.PTConfigs, error) {
	g := r.Distillation.PTHandler.GetSensors()

	configs := make([]*distillationproto.PTConfig, len(g))
	for i, elem := range g {
		configs[i] = ptConfigToRPC(&elem)
	}
	return &distillationproto.PTConfigs{Configs: configs}, nil
}

func (r *RPC) PTConfigure(ctx context.Context, config *distillationproto.PTConfig) (*distillationproto.PTConfig, error) {
	logger.Debug("PTConfigure")
	cfg := rpcToPTConfig(config)
	newCfg, err := r.Distillation.PTHandler.Configure(cfg)
	if err != nil {
		return nil, err
	}
	return ptConfigToRPC(&newCfg), nil
}

func (r *RPC) PTGetTemperatures(ctx context.Context, e *empty.Empty) (*distillationproto.PTTemperatures, error) {
	t := r.Distillation.PTHandler.Temperatures()
	return ptTemperatureToRPC(t), nil
}

func (r *RPC) GetPhaseCount(ctx context.Context, e *empty.Empty) (*distillationproto.ProcessPhaseCount, error) {
	cfg := r.Distillation.Process.GetConfig()
	s := &distillationproto.ProcessPhaseCount{Count: int32(cfg.PhaseNumber)}
	return s, nil
}

func (r *RPC) GetPhaseConfig(ctx context.Context, number *distillationproto.PhaseNumber) (*distillationproto.ProcessPhaseConfig, error) {
	logger.Debug("GetPhaseConfig")
	cfg := r.Distillation.Process.GetConfig()
	if uint(number.Number) >= cfg.PhaseNumber {
		return nil, errors.New("wrong phase number")
	}

	c := ProcessPhaseConfig{PhaseConfig: cfg.Phases[number.Number]}
	return processPhaseConfigToRpc(int(number.Number), c), nil
}

func (r *RPC) ConfigurePhaseCount(ctx context.Context, count *distillationproto.ProcessPhaseCount) (*distillationproto.ProcessPhaseCount, error) {
	r.Distillation.Process.SetPhaseNumber(uint(count.Count))
	return count, nil
}

func (r *RPC) ConfigurePhase(ctx context.Context, config *distillationproto.ProcessPhaseConfig) (*distillationproto.ProcessPhaseConfig, error) {
	conf := rpcToProcessPhaseConfig(config)
	if err := r.Distillation.configurePhase(uint(config.Number.Number), conf); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *RPC) ValidateConfig(ctx context.Context, e *empty.Empty) (*distillationproto.ProcessConfigValidation, error) {
	v := r.Distillation.ValidateConfig()
	return &distillationproto.ProcessConfigValidation{Valid: v.Valid, Error: v.Error}, nil
}

func (r *RPC) ConfigureProcess(ctx context.Context, config *distillationproto.ProcessConfig) (*distillationproto.ProcessConfig, error) {
	conf := rpcToProcessConfig(config)
	if err := r.Distillation.ConfigureProcess(conf); err != nil {
		return nil, err
	}
	return config, nil
}

func (r *RPC) Status(ctx context.Context, e *empty.Empty) (*distillationproto.ProcessStatus, error) {
	status := r.Distillation.Status()
	return processStatusToRPC(status), nil
}

func (r *RPC) EnableProcess(_ context.Context, c *distillationproto.ProcessConfig) (*distillationproto.ProcessConfig, error) {
	cfg := rpcToProcessConfig(c)
	err := r.Distillation.ConfigureProcess(cfg)
	if err != nil {
		return nil, err
	}
	return processConfigToRpc(cfg), nil
}

func (r *RPC) GetGlobalConfig(context.Context, *empty.Empty) (*distillationproto.ProcessGlobalConfig, error) {
	cfg := r.Distillation.Process.GetConfig()
	globalConfig := &distillationproto.ProcessGlobalConfig{
		Count:            int32(cfg.PhaseNumber),
		PhaseConfig:      make([]*distillationproto.ProcessPhaseConfig, len(cfg.Phases)),
		GlobalGPIOConfig: make([]*distillationproto.GPIOPhaseConfig, len(cfg.GlobalGPIO)),
		Sensors:          cfg.Sensors,
	}

	for i, gpio := range cfg.GlobalGPIO {
		globalConfig.GlobalGPIOConfig[i] = &distillationproto.GPIOPhaseConfig{
			ID:         gpio.ID,
			SensorID:   gpio.SensorID,
			TLow:       float32(gpio.TLow),
			THigh:      float32(gpio.THigh),
			Hysteresis: float32(gpio.Hysteresis),
			Inverted:   gpio.Inverted,
			Enabled:    gpio.Enabled,
		}
	}

	for i, conf := range cfg.Phases {
		globalConfig.PhaseConfig[i] = processPhaseConfigToRpc(i, ProcessPhaseConfig{conf})
	}

	globalConfig.Sensors = cfg.Sensors
	return globalConfig, nil
}

func (r *RPC) ConfigureGlobalGPIO(ctx context.Context, c *distillationproto.GlobalGPIOConfig) (*distillationproto.GlobalGPIOConfig, error) {
	config := make([]process.GPIOConfig, len(c.Configs))
	for i, conf := range c.Configs {
		config[i] = rpcToGPIOPhaseConfig(conf)
	}

	if err := r.Distillation.Process.SetGPIOGlobalConfig(config); err != nil {
		return nil, err
	}

	newConf := r.Distillation.Process.GetConfig().GlobalGPIO
	proto := new(distillationproto.GlobalGPIOConfig)
	proto.Configs = make([]*distillationproto.GPIOPhaseConfig, len(newConf))
	for i, conf := range newConf {
		proto.Configs[i] = gpioPhaseConfigToRpc(conf)
	}

	return proto, nil
}
