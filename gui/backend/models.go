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

package backend

import (
	"distillation/pkg/distillation"
	"distillation/pkg/process"
	"embedded/pkg/gpio"
	"gui/backend/parameters"
)

// Models allows us to create models.ts in frontend with needed structures
type Models struct{}

func (*Models) Temperature() parameters.Temperature {
	return parameters.Temperature{}
}

func (*Models) HeaterPhaseConfig() process.HeaterPhaseConfig {
	return process.HeaterPhaseConfig{}
}

func (*Models) GPIOPhaseConfig() process.GPIOConfig {
	return process.GPIOConfig{}
}

func (*Models) ProcessConfigValidation() distillation.ProcessConfigValidation {
	return distillation.ProcessConfigValidation{}
}

func (*Models) DistillationProcessStatus() distillation.ProcessStatus {
	return distillation.ProcessStatus{}
}

func (*Models) HeaterPhaseStatus() process.HeaterPhaseStatus {
	return process.HeaterPhaseStatus{}
}

func (*Models) TemperaturePhaseStatus() process.TemperaturePhaseStatus {
	return process.TemperaturePhaseStatus{}
}

func (*Models) GPIOPhaseStatus() process.GPIOPhaseStatus {
	return process.GPIOPhaseStatus{}
}

func (*Models) MoveToNextConfig() process.MoveToNextConfig {
	return process.MoveToNextConfig{}
}

func (*Models) GPIOActiveLevel() gpio.ActiveLevel {
	return gpio.High
}

func (*Models) MoveToNextStatus() process.MoveToNextStatus {
	return process.MoveToNextStatus{}
}

func (*Models) MoveToNextStatusTemperature() process.MoveToNextStatusTemperature {
	return process.MoveToNextStatusTemperature{}
}

func (*Models) TemperatureErrorCodeWrongID() int {
	return int(distillation.ErrorCodeWrongID)
}

func (*Models) TemperatureErrorCodeEmptyBuffer() int {
	return int(distillation.ErrorCodeEmptyBuffer)
}

func (*Models) TemperatureErrorCodeInternal() int {
	return int(distillation.ErrorCodeInternal)
}

func (*Models) ProcessStatus() ProcessStatus {
	return ProcessStatus{}
}

func (*Models) Update() Update {
	return Update{}
}

func (*Models) UpdateNextState() UpdateNextState {
	return UpdateNextState{}
}

func (*Models) UpdateStateStatus() UpdateStateStatus {
	return UpdateStateStatus{}
}
