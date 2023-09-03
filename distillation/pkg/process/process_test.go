package process_test

import (
	"errors"
	"time"

	"distillation/pkg/process"
	"distillation/pkg/process/mocks"

	"github.com/golang/mock/gomock"
)

func (pcs *ProcessConfigSuite) TestHappyPath_ConfigureOnFly() {
	t := pcs.Require()
	ctrl := gomock.NewController(pcs.T())

	heaterMock := mocks.NewMockHeater(ctrl)
	heaterMock.EXPECT().ID().Return("h1").AnyTimes()
	heaters := []process.Heater{heaterMock}

	sensorMock := mocks.NewMockSensor(ctrl)
	sensorMock.EXPECT().ID().Return("s1").AnyTimes()
	sensors := []process.Sensor{sensorMock}

	outputMock := mocks.NewMockOutput(ctrl)
	outputMock.EXPECT().ID().Return("o1").AnyTimes()
	outputs := []process.Output{outputMock}

	clockMock := mocks.NewMockClock(ctrl)

	p := process.New(process.WithClock(clockMock))
	p.UpdateOutputs(outputs)
	p.UpdateHeaters(heaters)
	p.UpdateSensors(sensors)

	cfg := process.Config{
		PhaseNumber: 2,
		Phases: []process.PhaseConfig{
			{
				Next: process.MoveToNextConfig{
					Type:            process.ByTime,
					SensorID:        "",
					SensorThreshold: 0,
					TimeLeft:        100,
				},
				Heaters: []process.HeaterPhaseConfig{
					{
						ID:    "h1",
						Power: 13,
					},
				},
				GPIO: []process.GPIOConfig{
					{
						Enabled:    true,
						ID:         "o1",
						SensorID:   "s1",
						TLow:       10,
						THigh:      20,
						Hysteresis: 1,
						Inverted:   false,
					},
				},
			},
			{
				Next: process.MoveToNextConfig{
					Type:            process.ByTemperature,
					SensorID:        "s1",
					SensorThreshold: 75.0,
					TimeLeft:        10,
				},
				Heaters: []process.HeaterPhaseConfig{
					{
						ID:    "h1",
						Power: 40,
					},
				},
				GPIO: []process.GPIOConfig{
					{
						Enabled:    true,
						ID:         "o1",
						SensorID:   "s1",
						TLow:       50,
						THigh:      51,
						Hysteresis: 1,
						Inverted:   true,
					},
				},
			},
		},
		GlobalGPIO: []process.GPIOConfig{
			{
				Enabled:    false,
				ID:         "o1",
				SensorID:   "s1",
				TLow:       10,
				THigh:      20,
				Hysteresis: 1,
				Inverted:   false,
			},
		},
	}
	p.SetPhaseNumber(cfg.PhaseNumber)
	t.Nil(p.SetGPIOGlobalConfig(cfg.GlobalGPIO))
	for i, conf := range cfg.Phases {
		t.Nil(p.SetPhaseConfig(uint(i), conf))
	}
	retTime := int64(0)
	// We expect:
	clockMock.EXPECT().Unix().Return(retTime).Times(2)
	sensorMock.EXPECT().Temperature().Return(1.1, nil)
	// 3. Set Power on heater
	heaterMock.EXPECT().SetPower(13).Return(nil)
	// 4. Handle GPIO
	outputMock.EXPECT().Set(false).Return(nil)
	s, err := p.Run()
	t.Nil(err)
	expectedStatus := process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(retTime, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    100,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Second time, just to move time
	retTime = 5
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(1.1, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(false).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    95,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Now, lets modify config and see what happens
	phase := cfg.Phases[0]
	phase.GPIO[0].Inverted = true
	phase.Heaters[0].Power = 56
	phase.Next.TimeLeft = 500
	t.Nil(p.SetPhaseConfig(uint(0), phase))

	// We expect that GPIO should be set now - it is inverted
	// Power of heater should be 56
	retTime = 10
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(1.1, nil)
	heaterMock.EXPECT().SetPower(56).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    490,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 56,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Move to second phase
	retTime = 501
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(1.1, nil)
	heaterMock.EXPECT().SetPower(40).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 1,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 40,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Now, lets modify config - TemperatureHoldSeconds and see what happens
	phase = cfg.Phases[1]
	phase.Next.TimeLeft = 500
	t.Nil(p.SetPhaseConfig(1, phase))
	retTime = 600
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(77.1, nil)
	heaterMock.EXPECT().SetPower(40).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 1,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 500,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 40,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 77.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Now some time elapsed, and we are changing hold again
	phase = cfg.Phases[1]
	phase.Next.TimeLeft = 800
	t.Nil(p.SetPhaseConfig(1, phase))
	retTime = 700
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(77.1, nil)
	heaterMock.EXPECT().SetPower(40).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 1,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 700,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 40,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 77.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)
}

func (pcs *ProcessConfigSuite) TestHappyPath_VerifyGPIOHandlingSinglePhase() {
	t := pcs.Require()
	ctrl := gomock.NewController(pcs.T())

	heaterMock := mocks.NewMockHeater(ctrl)
	heaterMock.EXPECT().ID().Return("h1").AnyTimes()
	heaters := []process.Heater{heaterMock}

	sensorMock := mocks.NewMockSensor(ctrl)
	sensorMock.EXPECT().ID().Return("s1").AnyTimes()
	sensors := []process.Sensor{sensorMock}

	outputMock := mocks.NewMockOutput(ctrl)
	outputMock.EXPECT().ID().Return("o1").AnyTimes()
	outputs := []process.Output{outputMock}

	clockMock := mocks.NewMockClock(ctrl)

	p := process.New(process.WithClock(clockMock))
	p.UpdateHeaters(heaters)
	p.UpdateOutputs(outputs)
	p.UpdateSensors(sensors)

	cfg := process.Config{
		PhaseNumber: 1,
		Phases: []process.PhaseConfig{
			{
				Next: process.MoveToNextConfig{
					Type:            process.ByTime,
					SensorID:        "",
					SensorThreshold: 0,
					TimeLeft:        100,
				},
				Heaters: []process.HeaterPhaseConfig{
					{
						ID:    "h1",
						Power: 13,
					},
				},
				GPIO: []process.GPIOConfig{
					{
						Enabled:    true,
						ID:         "o1",
						SensorID:   "s1",
						TLow:       10,
						THigh:      20,
						Hysteresis: 1,
						Inverted:   false,
					},
				},
			},
		},
		GlobalGPIO: []process.GPIOConfig{
			{
				Enabled:    false,
				ID:         "o1",
				SensorID:   "s1",
				TLow:       10,
				THigh:      20,
				Hysteresis: 1,
				Inverted:   false,
			},
		},
	}

	p.SetPhaseNumber(cfg.PhaseNumber)
	t.Nil(p.SetGPIOGlobalConfig(cfg.GlobalGPIO))
	for i, conf := range cfg.Phases {
		t.Nil(p.SetPhaseConfig(uint(i), conf))
	}
	retTime := int64(0)
	// We expect:
	// 1. Single call to Unix()
	clockMock.EXPECT().Unix().Return(retTime).Times(2)
	// 2. Twice Get temperature from sensor, so it can be placed in status
	sensorMock.EXPECT().Temperature().Return(1.1, nil)
	// 3. Set Power on heater
	heaterMock.EXPECT().SetPower(13).Return(nil)
	// 4. Handle GPIO
	outputMock.EXPECT().Set(false).Return(nil)
	s, err := p.Run()
	t.Nil(err)
	expectedStatus := process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(retTime, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    100,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Second: gpio should be set
	retTime = 5
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(15.1, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    95,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 15.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Third call, gpio should still be set
	retTime = 17
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(21.0, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    83,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 21,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Fourth call, gpio should still be set
	retTime = 17
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(9.0, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    83,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 9,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Fifth call, gpio should be off
	retTime = 17
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(8.9, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(false).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    83,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 8.9,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Sixth call, gpio still off
	retTime = 17
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(9.9, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(false).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    83,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 9.9,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Seventh call, gpio ON
	retTime = 17
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(10.0, nil)
	heaterMock.EXPECT().SetPower(13).Return(nil)
	outputMock.EXPECT().Set(true).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    83,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 10.0,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: true,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// FourthCall - finished
	retTime = 101
	clockMock.EXPECT().Unix().Return(retTime)
	sensorMock.EXPECT().Temperature().Return(100.1, nil)
	// Disable heater
	heaterMock.EXPECT().SetPower(0).Return(nil)
	outputMock.EXPECT().Set(false).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     false,
		Done:        true,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Unix(retTime, 0),
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    0,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 0,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 100.1,
			},
		},
		GPIO: []process.GPIOPhaseStatus{
			{
				ID:    "o1",
				State: false,
			},
		},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)
}

func (pcs *ProcessConfigSuite) TestHappyPath_SinglePhaseByTemperature() {
	t := pcs.Require()
	ctrl := gomock.NewController(pcs.T())

	clockMock := mocks.NewMockClock(ctrl)
	p := process.New(process.WithClock(clockMock))

	heater := mocks.NewMockHeater(ctrl)
	heater.EXPECT().ID().Return("h1").AnyTimes()
	p.UpdateHeaters([]process.Heater{heater})

	sensor := mocks.NewMockSensor(ctrl)
	sensor.EXPECT().ID().Return("s1").AnyTimes()
	p.UpdateSensors([]process.Sensor{sensor})

	cfg := process.Config{
		PhaseNumber: 1,
		Phases: []process.PhaseConfig{
			{
				Next: process.MoveToNextConfig{
					Type:            process.ByTemperature,
					SensorID:        "s1",
					SensorThreshold: 75.0,
					TimeLeft:        10,
				},
				Heaters: []process.HeaterPhaseConfig{
					{
						ID:    "h1",
						Power: 13,
					},
				},
				GPIO: nil,
			},
		},
	}

	p.SetPhaseNumber(cfg.PhaseNumber)
	t.Nil(p.SetGPIOGlobalConfig(cfg.GlobalGPIO))
	for i, conf := range cfg.Phases {
		t.Nil(p.SetPhaseConfig(uint(i), conf))
	}
	retTime := int64(0)

	// We expect:
	clockMock.EXPECT().Unix().Return(retTime).Times(2)
	// 2. Three times get temperature from sensor: built first Status, next() from byTemperature and built status again
	sensor.EXPECT().Temperature().Return(1.1, nil)
	// 3. Set Power on heater
	heater.EXPECT().SetPower(13).Return(nil)
	s, err := p.Run()
	t.Nil(err)
	expectedStatus := process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(retTime, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// On second Process we just expect single calls, except Temperature, which is called via ByTemperature
	retTime = 5
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(74.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 74.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Third call, temperature over threshold
	retTime = 99
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(75.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 75.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// FourthCall - temperature over threshold, time not elapsed
	retTime = 101
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(100.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 8,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 100.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Fifth - temperature under threshold, timeleft should be back to 10
	retTime = 103
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(74.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 74.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Sixth - temperature over threshold
	retTime = 150
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(75.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 10,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 75.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Seventh - temperature over threshold, time elapsed
	retTime = 160
	clockMock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(75.1, nil)
	// Disable heater
	heater.EXPECT().SetPower(0).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     false,
		Done:        true,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Unix(retTime, 0),
		Next: process.MoveToNextStatus{
			Type:     process.ByTemperature,
			TimeLeft: 0,
			Temperature: process.MoveToNextStatusTemperature{
				SensorID:        "s1",
				SensorThreshold: 75.0,
			},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 0,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 75.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)
}

func (pcs *ProcessConfigSuite) TestHappyPath_SinglePhaseByTime() {
	t := pcs.Require()
	ctrl := gomock.NewController(pcs.T())

	clock := mocks.NewMockClock(ctrl)
	p := process.New(process.WithClock(clock))

	heater := mocks.NewMockHeater(ctrl)
	heater.EXPECT().ID().Return("h1").AnyTimes()
	p.UpdateHeaters([]process.Heater{heater})

	sensor := mocks.NewMockSensor(ctrl)
	sensor.EXPECT().ID().Return("s1").AnyTimes()
	p.UpdateSensors([]process.Sensor{sensor})

	cfg := process.Config{
		PhaseNumber: 1,
		Phases: []process.PhaseConfig{
			{
				Next: process.MoveToNextConfig{
					Type:            process.ByTime,
					SensorID:        "",
					SensorThreshold: 0,
					TimeLeft:        100,
				},
				Heaters: []process.HeaterPhaseConfig{
					{
						ID:    "h1",
						Power: 13,
					},
				},
				GPIO: nil,
			},
		},
	}
	p.SetPhaseNumber(cfg.PhaseNumber)
	t.Nil(p.SetGPIOGlobalConfig(cfg.GlobalGPIO))
	for i, conf := range cfg.Phases {
		t.Nil(p.SetPhaseConfig(uint(i), conf))
	}

	retTime := int64(0)
	// We expect:
	// 1. Twice call to Unix()
	clock.EXPECT().Unix().Return(retTime).Times(2)
	// 2. Twice Get temperature from sensor, so it can be placed in status
	sensor.EXPECT().Temperature().Return(1.1, nil)
	// 3. Set Power on heater
	heater.EXPECT().SetPower(13).Return(nil)
	s, err := p.Run()
	t.Nil(err)
	expectedStatus := process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(retTime, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    100,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 1.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// On second Process we just expect single calls
	retTime = 5
	clock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(100.1, nil)
	heater.EXPECT().SetPower(13).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    95,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 13,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 100.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)

	// Third call, err on SetPower
	retTime = 99
	pwrErr := errors.New("hello there")
	clock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(150.1, nil)
	heater.EXPECT().SetPower(13).Return(pwrErr)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     true,
		Done:        false,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Time{},
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    1,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			// Empty, because error happened
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 150.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	// Check without errors,
	errs := s.Errors
	s.Errors = nil
	expectedStatus.Errors = nil
	t.EqualValues(expectedStatus, s)
	t.Len(errs, 1)
	t.Contains(errs[0], pwrErr.Error())

	// FourthCall - finished
	retTime = 101
	clock.EXPECT().Unix().Return(retTime)
	sensor.EXPECT().Temperature().Return(100.1, nil)
	// Disable heater
	heater.EXPECT().SetPower(0).Return(nil)
	s, err = p.Process()
	t.Nil(err)
	expectedStatus = process.Status{
		Running:     false,
		Done:        true,
		PhaseNumber: 0,
		StartTime:   time.Unix(0, 0),
		EndTime:     time.Unix(retTime, 0),
		Next: process.MoveToNextStatus{
			Type:        process.ByTime,
			TimeLeft:    0,
			Temperature: process.MoveToNextStatusTemperature{},
		},
		Heaters: []process.HeaterPhaseStatus{
			{
				process.HeaterPhaseConfig{
					ID:    "h1",
					Power: 0,
				},
			},
		},
		Temperature: []process.TemperaturePhaseStatus{
			{
				ID:          "s1",
				Temperature: 100.1,
			},
		},
		GPIO:   []process.GPIOPhaseStatus{},
		Errors: nil,
	}
	t.EqualValues(expectedStatus, s)
}
