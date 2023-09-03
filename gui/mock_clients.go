// MIT License
//
// Copyright (c) 2023 a-clap
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"distillation/pkg/distillation"
	"embedded/pkg/ds18b20"
	"embedded/pkg/embedded"
	"embedded/pkg/gpio"
	"embedded/pkg/max31865"
	"log"
	"osservice/pkg/wifi"
	"time"

	"gui/backend"
	"gui/backendmock"
)

func mockClients() []backend.Option {
	// HeaterClient - Mock
	heaterClient := backendmock.HeaterClient{}
	heaterClient.Globals = append(heaterClient.Globals,
		distillation.HeaterConfigGlobal{ID: "heater_1", Enabled: false},
		distillation.HeaterConfigGlobal{ID: "heater_2", Enabled: true},
		distillation.HeaterConfigGlobal{ID: "heater_3", Enabled: false},
		distillation.HeaterConfigGlobal{ID: "heater_4", Enabled: false},
	)

	dsClient := backendmock.DSClient{}
	dsClient.DS = append(dsClient.DS,
		distillation.DSConfig{DSSensorConfig: embedded.DSSensorConfig{SensorConfig: ds18b20.SensorConfig{
			Name:         "DS_1",
			ID:           "1",
			Resolution:   ds18b20.Resolution9Bit,
			Correction:   0,
			PollInterval: 500,
			Samples:      1,
		}}},
		distillation.DSConfig{DSSensorConfig: embedded.DSSensorConfig{SensorConfig: ds18b20.SensorConfig{
			Name:         "DS_2",
			ID:           "2",
			Resolution:   ds18b20.Resolution10Bit,
			Correction:   0,
			PollInterval: 1,
			Samples:      2,
		}}},
		distillation.DSConfig{DSSensorConfig: embedded.DSSensorConfig{SensorConfig: ds18b20.SensorConfig{
			Name:         "DS_3",
			ID:           "3",
			Resolution:   ds18b20.Resolution11Bit,
			Correction:   0,
			PollInterval: 1,
			Samples:      3,
		}}},
		distillation.DSConfig{DSSensorConfig: embedded.DSSensorConfig{SensorConfig: ds18b20.SensorConfig{
			Name:         "DS_4",
			ID:           "4",
			Resolution:   ds18b20.Resolution12Bit,
			Correction:   0,
			PollInterval: 1,
			Samples:      4,
		}}},
	)

	ptClient := backendmock.PTClient{}
	ptClient.PT = append(ptClient.PT,
		distillation.PTConfig{PTSensorConfig: embedded.PTSensorConfig{
			Enabled: false,
			SensorConfig: max31865.SensorConfig{
				Name:         "PT_1",
				ID:           "id_1",
				Correction:   0,
				ASyncPoll:    false,
				PollInterval: 1 * time.Second,
				Samples:      3,
			},
		}},
		distillation.PTConfig{PTSensorConfig: embedded.PTSensorConfig{
			Enabled: false,
			SensorConfig: max31865.SensorConfig{
				Name:         "PT_2",
				ID:           "id_2",
				Correction:   10.0,
				ASyncPoll:    true,
				PollInterval: 1 * time.Second,
				Samples:      7,
			},
		}},
		distillation.PTConfig{PTSensorConfig: embedded.PTSensorConfig{
			Enabled: false,
			SensorConfig: max31865.SensorConfig{
				Name:         "PT_3",
				ID:           "id_3",
				Correction:   12.0,
				ASyncPoll:    true,
				PollInterval: 1 * time.Second,
				Samples:      13,
			},
		}},
	)

	gpioClient := backendmock.GPIOClient{}
	gpioClient.GPIO = append(gpioClient.GPIO,
		distillation.GPIOConfig{GPIOConfig: embedded.GPIOConfig{Config: gpio.Config{
			ID:          "gpio_1",
			Direction:   gpio.DirOutput,
			ActiveLevel: gpio.Low,
			Value:       false,
		}}},
		distillation.GPIOConfig{GPIOConfig: embedded.GPIOConfig{Config: gpio.Config{
			ID:          "gpio_2",
			Direction:   gpio.DirOutput,
			ActiveLevel: gpio.High,
			Value:       true,
		}}},
		distillation.GPIOConfig{GPIOConfig: embedded.GPIOConfig{Config: gpio.Config{
			ID:          "gpio_3",
			Direction:   gpio.DirOutput,
			ActiveLevel: gpio.Low,
			Value:       true,
		}}},
	)

	phaseClient := backendmock.PhasesClient{}
	phaseClient.Init(3)

	w, err := wifi.New()
	if err != nil {
		log.Fatalln(err)
	}

	saver := backendmock.NewSaver()

	return []backend.Option{
		backend.WithHeaterClient(&heaterClient),
		backend.WithDSClient(&dsClient),
		backend.WithPTClient(&ptClient),
		backend.WithGPIOClient(&gpioClient),
		backend.WithPhaseClient(&phaseClient),
		backend.WithWifi(w),
		backend.WithLoadSaver(saver),
		backend.WithNet(backendmock.NetMock{}),
		backend.WithTime(backendmock.TimeMock{}),
		backend.WithUpdate(&backendmock.Update{}),
	}
}
