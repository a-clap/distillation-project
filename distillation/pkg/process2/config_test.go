package process2_test

import (
	"testing"

	"github.com/a-clap/distillation/pkg/process2"
	"github.com/a-clap/distillation/pkg/process2/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
)

type ProcessConfigSuite struct {
	suite.Suite
}

func TestProcess(t *testing.T) {
	suite.Run(t, new(ProcessConfigSuite))
}

func (pcs *ProcessConfigSuite) TestSetPhases() {
	t := pcs.Require()
	p := process2.New()

	p.SetPhaseNumber(5)
	cfg := p.GetConfig()
	t.EqualValues(cfg.PhaseNumber, 5)
	t.Len(cfg.Phases, 5)

	// Make it bigger
	p.SetPhaseNumber(7)
	cfg = p.GetConfig()
	t.EqualValues(cfg.PhaseNumber, 7)
	t.Len(cfg.Phases, 7)

	// Make it smaller
	p.SetPhaseNumber(3)
	cfg = p.GetConfig()
	t.EqualValues(cfg.PhaseNumber, 3)
	t.Len(cfg.Phases, 3)

	// Make it zero
	p.SetPhaseNumber(0)
	cfg = p.GetConfig()
	t.EqualValues(cfg.PhaseNumber, 0)
	t.Len(cfg.Phases, 0)
}

func (pcs *ProcessConfigSuite) TestGPIOGlobalConfig() {
	t := pcs.Require()
	args := []struct {
		name         string
		sensorIDs    []string
		outputs      []string
		globalConfig []process2.GPIOConfig
		err          error
	}{
		{
			name:      "duplicated ID",
			sensorIDs: nil,
			outputs:   nil,
			globalConfig: []process2.GPIOConfig{
				{
					ID: "o1",
				},
				{
					ID: "o1",
				},
			},
			err: process2.ErrDuplicatedID,
		},
		{
			name:      "no such GPIO",
			sensorIDs: nil,
			outputs:   nil,
			globalConfig: []process2.GPIOConfig{
				{
					ID: "o1",
				},
			},
			err: process2.ErrWrongGpioID,
		},
		{
			name:      "no such sensor",
			sensorIDs: nil,
			outputs:   []string{"o1"},
			globalConfig: []process2.GPIOConfig{
				{
					ID:       "o1",
					SensorID: "s1",
				},
			},
			err: process2.ErrWrongSensorID,
		},
		{
			name:      "all good",
			sensorIDs: []string{"s1"},
			outputs:   []string{"o1"},
			globalConfig: []process2.GPIOConfig{
				{
					ID:       "o1",
					SensorID: "s1",
				},
			},
			err: nil,
		},
	}
	for _, arg := range args {
		ctrl := gomock.NewController(pcs.T())

		pr := process2.New()
		if arg.sensorIDs != nil {
			sensors := make([]process2.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.outputs != nil {
			outputs := make([]process2.Output, len(arg.outputs))
			for i, elem := range arg.outputs {
				output := mocks.NewMockOutput(ctrl)
				output.EXPECT().ID().Return(elem).AnyTimes()
				outputs[i] = output
			}
			pr.UpdateOutputs(outputs)
		}

		err := pr.SetGPIOGlobalConfig(arg.globalConfig)
		if arg.err != nil {
			t.NotNil(err, arg.name)
			continue
		}
		t.Nil(err, arg.name)

		t.ElementsMatch(arg.globalConfig, pr.GetConfig().GlobalGPIO, arg.name)

	}

}
func (pcs *ProcessConfigSuite) TestConfigReflectsComponents() {
	args := []struct {
		name      string
		sensorIDs []string
		heaterIDs []string
		outputIDs []string
	}{
		{
			name:      "gpio",
			sensorIDs: nil,
			heaterIDs: nil,
			outputIDs: []string{"o1", "o2", "o3"},
		},
		{
			name:      "heaters",
			sensorIDs: nil,
			heaterIDs: []string{"H1", "H2"},
			outputIDs: nil,
		},
		{
			name:      "sensors",
			sensorIDs: []string{"s11", "s22"},
			heaterIDs: nil,
			outputIDs: nil,
		},
	}
	t := pcs.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(pcs.T())
		pr := process2.New()

		if arg.sensorIDs != nil {
			sensors := make([]process2.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.heaterIDs != nil {
			heaters := make([]process2.Heater, len(arg.heaterIDs))
			for i, elem := range arg.heaterIDs {
				heater := mocks.NewMockHeater(ctrl)
				heater.EXPECT().ID().Return(elem).AnyTimes()
				heaters[i] = heater
			}
			pr.UpdateHeaters(heaters)
		}

		if arg.outputIDs != nil {
			outputs := make([]process2.Output, len(arg.outputIDs))
			for i, elem := range arg.outputIDs {
				output := mocks.NewMockOutput(ctrl)
				output.EXPECT().ID().Return(elem).AnyTimes()
				outputs[i] = output
			}
			pr.UpdateOutputs(outputs)
		}

		cfg := pr.GetConfig()

		// Global GPIO
		t.Len(cfg.GlobalGPIO, len(arg.outputIDs), arg.name)
		for _, id := range arg.outputIDs {
			got := slices.ContainsFunc(cfg.GlobalGPIO, func(config process2.GPIOConfig) bool {
				return config.ID == id
			})
			t.True(got, "should contain ID: "+id)
		}

		for i, elem := range cfg.Phases {
			// GPIO
			t.Len(elem.GPIO, len(arg.outputIDs), arg.name, i)
			for _, id := range arg.outputIDs {
				got := slices.ContainsFunc(elem.GPIO, func(config process2.GPIOConfig) bool {
					return config.ID == id
				})
				t.True(got, "should contain ID: "+id)
			}
			// Heaters
			t.Len(elem.Heaters, len(arg.heaterIDs), arg.name, i)
			for _, id := range arg.heaterIDs {
				got := slices.ContainsFunc(elem.Heaters, func(config process2.HeaterPhaseConfig) bool {
					return config.ID == id
				})
				t.True(got, "should contain ID: "+id)
			}

			// Sensors
			t.Len(elem.Next.Sensors, len(arg.sensorIDs), arg.name, i)
		}
	}
}

func (pcs *ProcessConfigSuite) TestConfigurePhase_HeatersError() {
	t := pcs.Require()

	args := []struct {
		name          string
		heaterIDs     []string
		heatersConfig []process2.HeaterPhaseConfig
		err           error
	}{
		{
			name:      "wrong id of heater",
			heaterIDs: []string{"h1"},
			heatersConfig: []process2.HeaterPhaseConfig{
				{
					ID:    "h2",
					Power: 13,
				},
			},
			err: process2.ErrWrongHeaterID,
		},
		{
			name:      "power of heater over 100",
			heaterIDs: []string{"h1"},
			heatersConfig: []process2.HeaterPhaseConfig{
				{
					ID:    "h1",
					Power: 101,
				},
			},
			err: process2.ErrWrongHeaterPower,
		},
		{
			name:          "lack of heater configuration",
			heaterIDs:     []string{"h1"},
			heatersConfig: nil,
			err:           process2.ErrHeaterConfigDiffersFromHeatersLen,
		},

		{
			name:      "all good",
			heaterIDs: []string{"h1"},
			heatersConfig: []process2.HeaterPhaseConfig{
				{
					ID:    "h1",
					Power: 15,
				},
			},
			err: nil,
		},
	}
	for _, arg := range args {
		// Always good config - except heaters
		phaseConfig := process2.PhaseConfig{
			Next: process2.MoveToNextConfig{
				Type:            process2.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        3,
			},
			Heaters: nil,
			GPIO:    nil,
		}
		ctrl := gomock.NewController(pcs.T())
		pr := process2.New()
		if arg.heaterIDs != nil {
			heaters := make([]process2.Heater, len(arg.heaterIDs))
			for i, elem := range arg.heaterIDs {
				heater := mocks.NewMockHeater(ctrl)
				heater.EXPECT().ID().Return(elem).AnyTimes()
				heaters[i] = heater
			}
			pr.UpdateHeaters(heaters)
		}

		s := mocks.NewMockSensor(ctrl)
		s.EXPECT().ID().Return("s1").AnyTimes()
		pr.UpdateSensors([]process2.Sensor{s})

		phaseConfig.Heaters = arg.heatersConfig

		pr.SetPhaseNumber(1)
		err := pr.SetPhaseConfig(0, phaseConfig)
		if arg.err != nil {
			t.NotNil(err, arg.name)
			t.ErrorContains(err, arg.err.Error(), arg.name)
			continue
		}
		t.Nil(err, arg.name)
		cfg := pr.GetConfig()
		t.EqualValues(phaseConfig, cfg.Phases[0], arg.name)
	}
}

func (pcs *ProcessConfigSuite) TestConfigurePhase_GPIOErrors() {
	t := pcs.Require()

	args := []struct {
		name       string
		outputIDs  []string
		sensorIDs  []string
		gpioConfig []process2.GPIOConfig
		err        error
	}{
		{
			name:      "wrong id of gpio",
			outputIDs: []string{"h2"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process2.GPIOConfig{
				{
					ID:       "h2",
					SensorID: "s2",
				},
			},
			err: process2.ErrWrongSensorID,
		},
		{
			name:      "wrong gpio ID",
			outputIDs: []string{"g1"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process2.GPIOConfig{
				{
					ID:       "g2",
					SensorID: "s1",
				},
			},
			err: process2.ErrWrongGpioID,
		},
		{
			name:       "lack of gpio config",
			outputIDs:  []string{"g1"},
			sensorIDs:  []string{"s1"},
			gpioConfig: nil,
			err:        process2.ErrDifferentGPIOSConfig,
		},

		{
			name:      "all good",
			outputIDs: []string{"h1"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process2.GPIOConfig{
				{
					ID:       "h1",
					SensorID: "s1",
				},
			},
			err: nil,
		},
	}
	for _, arg := range args {
		// Always good config - except heaters
		phaseConfig := process2.PhaseConfig{
			Next: process2.MoveToNextConfig{
				Type:            process2.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        3,
			},
			Heaters: []process2.HeaterPhaseConfig{
				{ID: "h1", Power: 13},
			},
			GPIO: nil,
		}
		ctrl := gomock.NewController(pcs.T())
		pr := process2.New()

		if arg.sensorIDs != nil {
			sensors := make([]process2.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.outputIDs != nil {
			outputs := make([]process2.Output, len(arg.outputIDs))
			for i, elem := range arg.outputIDs {
				output := mocks.NewMockOutput(ctrl)
				output.EXPECT().ID().Return(elem).AnyTimes()
				outputs[i] = output
			}
			pr.UpdateOutputs(outputs)
		}
		{
			heater := mocks.NewMockHeater(ctrl)
			heater.EXPECT().ID().Return("h1").AnyTimes()
			pr.UpdateHeaters([]process2.Heater{heater})
		}
		phaseConfig.GPIO = arg.gpioConfig
		pr.SetPhaseNumber(1)
		err := pr.SetPhaseConfig(0, phaseConfig)
		if arg.err != nil {
			t.NotNil(err, arg.name)
			t.ErrorContains(err, arg.err.Error(), arg.name)
			continue
		}
		t.Nil(err, arg.name)
		cfg := pr.GetConfig()
		t.EqualValues(phaseConfig, cfg.Phases[0], arg.name)

	}
}

func (pcs *ProcessConfigSuite) TestConfigurePhase_NextConfig() {
	t := pcs.Require()

	args := []struct {
		name             string
		moveToNextConfig process2.MoveToNextConfig
		sensorIDs        []string
		err              error
	}{
		{
			name: "byTime - time can't be 0",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process2.ErrByTimeWrongTime,
		},
		{
			name: "byTime - seconds under 0",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        -1,
			},
			sensorIDs: []string{"s1"},
			err:       process2.ErrByTimeWrongTime,
		},
		{
			name: "byTime - all good",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        1,
			},
			sensorIDs: []string{"s1"},
			err:       nil,
		},
		{
			name: "byTemperature - wrong sensor",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.ByTemperature,
				SensorID:        "s2",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process2.ErrByTemperatureWrongID,
		},
		{
			name: "byTemperature - weird type",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.MoveToNextType(3),
				SensorID:        "s1",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process2.ErrUnknownType,
		},
		{
			name: "byTemperature - all good, threshold/hold can be 0",
			moveToNextConfig: process2.MoveToNextConfig{
				Type:            process2.ByTemperature,
				SensorID:        "s1",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       nil,
		},
	}
	for _, arg := range args {
		// Always good config - except Next
		phaseConfig := process2.PhaseConfig{
			Next: process2.MoveToNextConfig{
				Type:            0,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			Heaters: []process2.HeaterPhaseConfig{
				{ID: "h1", Power: 13},
			},
			GPIO: nil,
		}
		pr := process2.New()
		ctrl := gomock.NewController(pcs.T())

		if arg.sensorIDs != nil {
			sensors := make([]process2.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}
		{
			heater := mocks.NewMockHeater(ctrl)
			heater.EXPECT().ID().Return("h1").AnyTimes()
			pr.UpdateHeaters([]process2.Heater{heater})
		}

		phaseConfig.Next = arg.moveToNextConfig
		pr.SetPhaseNumber(1)
		err := pr.SetPhaseConfig(0, phaseConfig)
		if arg.err != nil {
			t.NotNil(err, arg.name)
			t.ErrorContains(err, arg.err.Error(), arg.name)
			continue
		}
		t.Nil(err, arg.name)
		cfg := pr.GetConfig()
		t.EqualValues(phaseConfig, cfg.Phases[0], arg.name)
	}
}
