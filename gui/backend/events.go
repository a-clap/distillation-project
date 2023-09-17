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

// Events are just way to return constants in UI
type Events struct{}

const (
	NotifyHeaters           = "rcv:heaters"
	NotifyDSConfig          = "rcv:dscfg"
	NotifyDSTemperature     = "rcv:dstmp"
	NotifyPTConfig          = "rcv:ptcfg"
	NotifyPTTemperature     = "rcv:pttmp"
	NotifyGPIO              = "rcv:gpio"
	NotifyError             = "rcv:error"
	NotifyPhasesConfig      = "rcv:phase_config"
	NotifyPhasesValidate    = "rcv:phase_validate"
	NotifyPhasesPhaseConfig = "rcv:phase_phase_config"
	NotifyPhasesPhaseCount  = "rcv:phase_count"
	NotifyPhasesStatus      = "rcv:phase_status"
	NotifyGlobalConfig      = "rcv:global_config"
	NotifyUpdate            = "rcv:update"
	NotifyUpdateStatus      = "rcv:update_status"
	NotifyUpdateNextState   = "rcv:update_next_state"
)

func (Events) NotifyHeaters() string {
	return NotifyHeaters
}

func (Events) NotifyDSConfig() string {
	return NotifyDSConfig
}

func (Events) NotifyDSTemperature() string {
	return NotifyDSTemperature
}

func (Events) NotifyPTConfig() string {
	return NotifyPTConfig
}

func (Events) NotifyPTTemperature() string {
	return NotifyPTTemperature
}

func (Events) NotifyGPIO() string {
	return NotifyGPIO
}

func (Events) NotifyError() string {
	return NotifyError
}

func (Events) NotifyPhasesConfig() string {
	return NotifyPhasesConfig
}

func (Events) NotifyPhasesValidate() string {
	return NotifyPhasesValidate
}

func (Events) NotifyPhasesPhaseConfig() string {
	return NotifyPhasesPhaseConfig
}

func (Events) NotifyPhasesPhaseCount() string {
	return NotifyPhasesPhaseCount
}

func (Events) NotifyPhasesStatus() string {
	return NotifyPhasesStatus
}

func (Events) NotifyGlobalConfig() string {
	return NotifyGlobalConfig
}

func (Events) NotifyUpdate() string {
	return NotifyUpdate
}

func (Events) NotifyUpdateStatus() string {
	return NotifyUpdateStatus
}
func (Events) NotifyUpdateNextState() string {
	return NotifyUpdateNextState
}
