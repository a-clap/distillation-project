package process

//go:generate mockgen -destination=mocks/mocks.go -package=mocks . Clock,Heater,Sensor,Output

import (
	"fmt"
	"sync/atomic"
	"time"

	"golang.org/x/exp/slices"
)

type Process struct {
	*config
	running     atomic.Bool
	phaseConfig *PhaseConfig
	phaseEnd    endCondition
	status      Status
	stamp       int64
}

func New(opts ...Option) *Process {
	p := &Process{
		config:  newConfig(),
		running: atomic.Bool{},
	}
	for _, opt := range opts {
		opt(p)
	}

	p.SetPhaseNumber(3)
	return p
}

func (p *Process) Run() (Status, error) {
	if err := p.Validate(); err != nil {
		return Status{}, err
	}

	if !p.Running() {
		p.init()
	}

	return p.Process()
}

func (p *Process) Process() (Status, error) {
	p.status.Errors = nil
	// Update processing stamp
	p.stamp = p.clock.Unix()
	// Update temperatures in Status
	p.updateTemperatures()

	var end bool
	if end, p.status.Next.TimeLeft = p.phaseEnd.end(); end {
		p.moveToPhase(p.status.PhaseNumber + 1)
	}

	p.handleHeaters()
	p.handleGpio()

	return p.status, nil
}

func (p *Process) Running() bool {
	return p.running.Load()
}

func (p *Process) MoveToNext() (Status, error) {
	if !p.Running() {
		return Status{}, ErrNotRunning
	}
	p.moveToPhase(p.status.PhaseNumber + 1)
	return p.status, nil
}

func (p *Process) Finish() (Status, error) {
	if !p.Running() {
		return Status{}, ErrNotRunning
	}
	p.finish()
	return p.status, nil
}

func (p *Process) SetPhaseConfig(nb uint, conf PhaseConfig) error {
	updatingCurrentPhase := p.Running() && p.status.PhaseNumber == nb && p.phaseConfig != nil
	var cfg PhaseConfig
	var status Status
	if updatingCurrentPhase {
		cfg = *p.phaseConfig
		status = p.status
	}

	if err := p.config.SetPhaseConfig(nb, conf); err != nil {
		return err
	}
	if updatingCurrentPhase {
		p.updateNextCondition()
		// Changing between types forces starting phase from the beginning
		if cfg.Next.Type == p.phaseConfig.Next.Type {
			timeElapsed := cfg.Next.TimeLeft - status.Next.TimeLeft
			newTimeleft := p.phaseConfig.Next.TimeLeft - timeElapsed

			// If more time elapsed, then we will move to next phase on next iteration
			if newTimeleft <= 0 {
				newTimeleft = 0
			}

			if p.phaseConfig.Next.Type == ByTime {
				if b, ok := p.phaseEnd.(*endConditionTime); ok {
					b.leftTime = newTimeleft
				}
			} else {
				if b, ok := p.phaseEnd.(*endConditionTemperature); ok {
					b.endTime.leftTime = newTimeleft
				}
			}

		}

	}
	return nil
}

func (p *Process) updateNextCondition() {
	timeNow := func() int64 {
		return p.stamp
	}

	if p.phaseConfig.Next.Type == ByTime {
		p.phaseEnd = newEndConditionTime(p.phaseConfig.Next.TimeLeft, timeNow)
	} else {
		temperatureNow := func() float64 {
			idx := slices.IndexFunc(p.status.Temperature, func(status TemperaturePhaseStatus) bool {
				return status.ID == p.phaseConfig.Next.SensorID
			})
			if idx == -1 {
				// TODO: how to send error?
				return -300
			}
			return p.status.Temperature[idx].Temperature
		}

		p.phaseEnd = newEndConditionTemperature(p.phaseConfig.Next.TimeLeft, timeNow, p.phaseConfig.Next.SensorThreshold, temperatureNow)
		p.status.Next.Temperature = MoveToNextStatusTemperature{
			SensorID:        p.phaseConfig.Next.SensorID,
			SensorThreshold: p.phaseConfig.Next.SensorThreshold,
		}
	}
	p.status.Next.Type = p.phaseConfig.Next.Type
	// Update timeLeft
	_, p.status.Next.TimeLeft = p.phaseEnd.end()
}

func (p *Process) moveToPhase(next uint) {
	if next >= p.Config.PhaseNumber {
		p.finish()
		return
	}

	p.status.PhaseNumber = next
	p.phaseConfig = &p.config.Phases[next]
	p.updateNextCondition()
}

func (p *Process) updateTemperatures() {
	p.status.Temperature = make([]TemperaturePhaseStatus, 0, len(p.sensors))
	for id, sensor := range p.sensors {
		// TODO: how to handle error?
		tmp, err := sensor.Temperature()
		if err != nil {
			err = fmt.Errorf("%w: read Temperature failed on ID: %v", err, id)
			p.status.Errors = append(p.status.Errors, err.Error())
			continue
		}
		p.status.Temperature = append(p.status.Temperature, TemperaturePhaseStatus{
			ID:          id,
			Temperature: tmp,
		})
	}
}

func (p *Process) handleHeaters() {
	p.status.Heaters = make([]HeaterPhaseStatus, 0, len(p.phaseConfig.Heaters))
	for _, config := range p.phaseConfig.Heaters {
		heater := p.heaters[config.ID]
		pwr := 0

		if p.status.Running {
			pwr = config.Power
		}

		if err := heater.SetPower(pwr); err != nil {
			err = fmt.Errorf("%w: on ID: %v, SetPower: %v", err, config.ID, config.Power)
			p.status.Errors = append(p.status.Errors, err.Error())
			continue
		}

		p.status.Heaters = append(p.status.Heaters, HeaterPhaseStatus{HeaterPhaseConfig{
			ID:    config.ID,
			Power: pwr,
		}})
	}
}

func (p *Process) handleSingleGPIO(config GPIOConfig) {
	gpio, ok := p.outputs[config.ID]
	if !ok {
		err := fmt.Errorf("output with ID: %v not found", config.ID)
		p.status.Errors = append(p.status.Errors, err.Error())
		return
	}

	gpioValue := false
	if p.status.Running && config.Enabled {
		idx := slices.IndexFunc(p.status.Temperature, func(status TemperaturePhaseStatus) bool {
			return status.ID == config.SensorID
		})
		if idx == -1 {
			err := fmt.Errorf("sensor with ID: %v not found", config.SensorID)
			p.status.Errors = append(p.status.Errors, err.Error())
			return
		}
		t := p.status.Temperature[idx].Temperature

		if gpio.inRange {
			// Last time was in range
			gpio.inRange = t >= (config.TLow-config.Hysteresis) && t <= (config.THigh+config.Hysteresis)
		} else {
			// Out of range, need to hit tlow or thigh
			gpio.inRange = t >= config.TLow && t <= config.THigh
		}
		gpioValue = gpio.inRange

		if config.Inverted {
			gpioValue = !gpioValue
		}
	}

	if err := gpio.Output.Set(gpioValue); err != nil {
		err := fmt.Errorf("%w: on gpio ID: %v, on Set with value: %v", err, config.ID, gpioValue)
		p.status.Errors = append(p.status.Errors, err.Error())
		return
	}

	p.status.GPIO = append(p.status.GPIO, GPIOPhaseStatus{
		ID:    config.ID,
		State: gpioValue,
	})
}

func (p *Process) handleGpio() {
	// Create slice of GPIOConfigs which we should handle
	configs := make([]GPIOConfig, 0, len(p.GlobalGPIO))
	// If GPIO is configured in PhaseConfig, it has higher priority
	ids := make(map[string]bool, len(p.GlobalGPIO))
	for _, conf := range p.phaseConfig.GPIO {
		if conf.Enabled {
			ids[conf.ID] = true
			configs = append(configs, conf)
		}
	}
	// Now look at global configs
	for _, conf := range p.GlobalGPIO {
		if v, ok := ids[conf.ID]; v && ok {
			// 	We already configured this gpio
			continue
		}
		configs = append(configs, conf)
	}

	p.status.GPIO = make([]GPIOPhaseStatus, 0, len(configs))
	for _, config := range configs {
		p.handleSingleGPIO(config)
	}
}

func (p *Process) init() {
	p.running.Store(true)
	p.status.Running = true
	p.status.Done = false
	p.stamp = p.clock.Unix()
	p.status.StartTime = time.Unix(p.stamp, 0)
	p.status.EndTime = time.Time{}
	p.moveToPhase(0)
}

func (p *Process) finish() {
	p.status.Done = true
	p.status.Running = false
	p.status.EndTime = time.Unix(p.stamp, 0)
	p.running.Store(false)
}
