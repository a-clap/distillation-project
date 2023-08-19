package process_test

import (
	"testing"

	"distillation/pkg/process"
	"distillation/pkg/process/mocks"

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

func (pcs *ProcessConfigSuite) TestConfig_AlwaysReturnsCorrectConfig() {
	args := []struct {
		name      string
		sensorIDs []string
		heaterIDs []string
		outputIDs []string
	}{
		{
			name:      "minimum components",
			sensorIDs: nil,
			heaterIDs: []string{"h1"},
			outputIDs: nil,
		},
		{
			name:      "heater and sensors",
			sensorIDs: []string{"s1"},
			heaterIDs: []string{"h1"},
			outputIDs: nil,
		},
		{
			name:      "heater, sensors and outputs",
			sensorIDs: []string{"s1"},
			heaterIDs: []string{"h1"},
			outputIDs: []string{"o1"},
		},
	}
	for _, arg := range args {
		ctrl := gomock.NewController(pcs.T())

		t := pcs.Require()
		p := process.New()

		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			p.UpdateSensors(sensors)
		}

		if arg.outputIDs != nil {
			outputs := make([]process.Output, len(arg.outputIDs))
			for i, elem := range arg.outputIDs {
				output := mocks.NewMockOutput(ctrl)
				output.EXPECT().ID().Return(elem).AnyTimes()
				outputs[i] = output
			}
			p.UpdateOutputs(outputs)
		}

		if arg.heaterIDs != nil {
			heaters := make([]process.Heater, len(arg.heaterIDs))
			for i, elem := range arg.heaterIDs {
				heater := mocks.NewMockHeater(ctrl)
				heater.EXPECT().ID().Return(elem).AnyTimes()
				heaters[i] = heater
			}
			p.UpdateHeaters(heaters)
		}

		cfg := p.GetConfig()
		t.Nil(p.SetGPIOGlobalConfig(cfg.GlobalGPIO))
		for i, conf := range cfg.Phases {
			t.Nil(p.SetPhaseConfig(uint(i), conf))
		}
	}

}
func (pcs *ProcessConfigSuite) TestConfigAvailableSensors() {
	args := []struct {
		name      string
		sensorIDs []string
	}{
		{
			name:      "single sensor",
			sensorIDs: []string{"s1"},
		},
		{
			name:      "two sensors",
			sensorIDs: []string{"s1", "s2"},
		},
		{
			name:      "many",
			sensorIDs: []string{"s1", "s2", "s5", "s7", "s12"},
		},
		{
			name:      "empty",
			sensorIDs: []string{},
		},
	}
	t := pcs.Require()
	for _, arg := range args {
		ctrl := gomock.NewController(pcs.T())

		pr := process.New()
		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}
		cfg := pr.GetConfig()
		t.ElementsMatch(cfg.Sensors, arg.sensorIDs)

	}

}

func (pcs *ProcessConfigSuite) TestSetPhases() {
	t := pcs.Require()
	p := process.New()

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
		globalConfig []process.GPIOConfig
		err          error
	}{
		{
			name:      "duplicated ID",
			sensorIDs: nil,
			outputs:   nil,
			globalConfig: []process.GPIOConfig{
				{
					ID: "o1",
				},
				{
					ID: "o1",
				},
			},
			err: process.ErrDuplicatedID,
		},
		{
			name:      "no such GPIO",
			sensorIDs: nil,
			outputs:   nil,
			globalConfig: []process.GPIOConfig{
				{
					ID: "o1",
				},
			},
			err: process.ErrWrongGpioID,
		},
		{
			name:      "no such sensor",
			sensorIDs: nil,
			outputs:   []string{"o1"},
			globalConfig: []process.GPIOConfig{
				{
					ID:       "o1",
					SensorID: "s1",
				},
			},
			err: process.ErrWrongSensorID,
		},
		{
			name:      "all good",
			sensorIDs: []string{"s1"},
			outputs:   []string{"o1"},
			globalConfig: []process.GPIOConfig{
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

		pr := process.New()
		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.outputs != nil {
			outputs := make([]process.Output, len(arg.outputs))
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
		pr := process.New()

		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.heaterIDs != nil {
			heaters := make([]process.Heater, len(arg.heaterIDs))
			for i, elem := range arg.heaterIDs {
				heater := mocks.NewMockHeater(ctrl)
				heater.EXPECT().ID().Return(elem).AnyTimes()
				heaters[i] = heater
			}
			pr.UpdateHeaters(heaters)
		}

		if arg.outputIDs != nil {
			outputs := make([]process.Output, len(arg.outputIDs))
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
			got := slices.ContainsFunc(cfg.GlobalGPIO, func(config process.GPIOConfig) bool {
				return config.ID == id
			})
			t.True(got, "should contain ID: "+id)
		}

		for i, elem := range cfg.Phases {
			// GPIO
			t.Len(elem.GPIO, len(arg.outputIDs), arg.name, i)
			for _, id := range arg.outputIDs {
				got := slices.ContainsFunc(elem.GPIO, func(config process.GPIOConfig) bool {
					return config.ID == id
				})
				t.True(got, "should contain ID: "+id)
			}
			// Heaters
			t.Len(elem.Heaters, len(arg.heaterIDs), arg.name, i)
			for _, id := range arg.heaterIDs {
				got := slices.ContainsFunc(elem.Heaters, func(config process.HeaterPhaseConfig) bool {
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
		heatersConfig []process.HeaterPhaseConfig
		err           error
	}{
		{
			name:      "wrong id of heater",
			heaterIDs: []string{"h1"},
			heatersConfig: []process.HeaterPhaseConfig{
				{
					ID:    "h2",
					Power: 13,
				},
			},
			err: process.ErrWrongHeaterID,
		},
		{
			name:      "power of heater over 100",
			heaterIDs: []string{"h1"},
			heatersConfig: []process.HeaterPhaseConfig{
				{
					ID:    "h1",
					Power: 101,
				},
			},
			err: process.ErrWrongHeaterPower,
		},
		{
			name:          "lack of heater configuration",
			heaterIDs:     []string{"h1"},
			heatersConfig: nil,
			err:           process.ErrHeaterConfigDiffersFromHeatersLen,
		},

		{
			name:      "all good",
			heaterIDs: []string{"h1"},
			heatersConfig: []process.HeaterPhaseConfig{
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
		phaseConfig := process.PhaseConfig{
			Next: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        3,
			},
			Heaters: nil,
			GPIO:    nil,
		}
		ctrl := gomock.NewController(pcs.T())
		pr := process.New()
		if arg.heaterIDs != nil {
			heaters := make([]process.Heater, len(arg.heaterIDs))
			for i, elem := range arg.heaterIDs {
				heater := mocks.NewMockHeater(ctrl)
				heater.EXPECT().ID().Return(elem).AnyTimes()
				heaters[i] = heater
			}
			pr.UpdateHeaters(heaters)
		}

		s := mocks.NewMockSensor(ctrl)
		s.EXPECT().ID().Return("s1").AnyTimes()
		pr.UpdateSensors([]process.Sensor{s})

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
		gpioConfig []process.GPIOConfig
		err        error
	}{
		{
			name:      "wrong id of gpio",
			outputIDs: []string{"h2"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process.GPIOConfig{
				{
					ID:       "h2",
					SensorID: "s2",
				},
			},
			err: process.ErrWrongSensorID,
		},
		{
			name:      "wrong gpio ID",
			outputIDs: []string{"g1"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process.GPIOConfig{
				{
					ID:       "g2",
					SensorID: "s1",
				},
			},
			err: process.ErrWrongGpioID,
		},
		{
			name:       "lack of gpio config",
			outputIDs:  []string{"g1"},
			sensorIDs:  []string{"s1"},
			gpioConfig: nil,
			err:        process.ErrDifferentGPIOSConfig,
		},

		{
			name:      "all good",
			outputIDs: []string{"h1"},
			sensorIDs: []string{"s1"},
			gpioConfig: []process.GPIOConfig{
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
		phaseConfig := process.PhaseConfig{
			Next: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        3,
			},
			Heaters: []process.HeaterPhaseConfig{
				{ID: "h1", Power: 13},
			},
			GPIO: nil,
		}
		ctrl := gomock.NewController(pcs.T())
		pr := process.New()

		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
			for i, elem := range arg.sensorIDs {
				sensor := mocks.NewMockSensor(ctrl)
				sensor.EXPECT().ID().Return(elem).AnyTimes()
				sensors[i] = sensor
			}
			pr.UpdateSensors(sensors)
		}

		if arg.outputIDs != nil {
			outputs := make([]process.Output, len(arg.outputIDs))
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
			pr.UpdateHeaters([]process.Heater{heater})
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
		moveToNextConfig process.MoveToNextConfig
		sensorIDs        []string
		err              error
	}{
		{
			name: "byTime - time can't be 0",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process.ErrByTimeWrongTime,
		},
		{
			name: "byTime - seconds under 0",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        -1,
			},
			sensorIDs: []string{"s1"},
			err:       process.ErrByTimeWrongTime,
		},
		{
			name: "byTime - all good",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        1,
			},
			sensorIDs: []string{"s1"},
			err:       nil,
		},
		{
			name: "byTemperature - wrong sensor",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.ByTemperature,
				SensorID:        "s2",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process.ErrByTemperatureWrongID,
		},
		{
			name: "byTemperature - weird type",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.MoveToNextType(3),
				SensorID:        "s1",
				SensorThreshold: 0,
				TimeLeft:        0,
			},
			sensorIDs: []string{"s1"},
			err:       process.ErrUnknownType,
		},
		{
			name: "byTemperature - all good, threshold/hold can be 0",
			moveToNextConfig: process.MoveToNextConfig{
				Type:            process.ByTemperature,
				SensorID:        "s1",
				SensorThreshold: 0,
				TimeLeft:        1,
			},
			sensorIDs: []string{"s1"},
			err:       nil,
		},
	}
	for _, arg := range args {
		// Always good config - except Next
		phaseConfig := process.PhaseConfig{
			Next: process.MoveToNextConfig{
				Type:            process.ByTime,
				SensorID:        "",
				SensorThreshold: 0,
				TimeLeft:        1,
			},
			Heaters: []process.HeaterPhaseConfig{
				{ID: "h1", Power: 13},
			},
			GPIO: nil,
		}
		pr := process.New()
		ctrl := gomock.NewController(pcs.T())

		if arg.sensorIDs != nil {
			sensors := make([]process.Sensor, len(arg.sensorIDs))
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
			pr.UpdateHeaters([]process.Heater{heater})
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
