/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package process2

import (
	"log"
	"time"

	"github.com/samber/lo"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

type MoveToNextType int

const (
	ByTime MoveToNextType = iota
	ByTemperature
)

type Config struct {
	PhaseNumber uint          `json:"phase_number"`
	Phases      []PhaseConfig `json:"phases"`
	GlobalGPIO  []GPIOConfig  `json:"global_gpio"`
}

type PhaseConfig struct {
	Next    MoveToNextConfig    `json:"next"`
	Heaters []HeaterPhaseConfig `json:"heaters"`
	GPIO    []GPIOConfig        `json:"gpio"`
}

type MoveToNextConfig struct {
	Type            MoveToNextType `json:"type"`
	Sensors         []string       `json:"sensors"`
	SensorID        string         `json:"sensor_id"`
	SensorThreshold float64        `json:"sensor_threshold"`
	TimeLeft        int64          `json:"time_left"`
}

type HeaterPhaseConfig struct {
	ID    string `json:"ID"`
	Power int    `json:"power"`
}

type GPIOConfig struct {
	Enabled    bool    `json:"enabled"`
	ID         string  `json:"id"`
	SensorID   string  `json:"sensor_id"`
	TLow       float64 `json:"t_low"`
	THigh      float64 `json:"t_high"`
	Hysteresis float64 `json:"hysteresis"`
	Inverted   bool    `json:"inverted"`
}

type HeaterPhaseStatus struct {
	HeaterPhaseConfig
}

type TemperaturePhaseStatus struct {
	ID          string  `json:"ID"`
	Temperature float64 `json:"temperature"`
}

type GPIOPhaseStatus struct {
	ID    string `json:"id"`
	State bool   `json:"state"`
}

type MoveToNextStatusTemperature struct {
	SensorID        string  `json:"sensor_id"`
	SensorThreshold float64 `json:"sensor_threshold"`
}

type MoveToNextStatus struct {
	Type        MoveToNextType              `json:"type"`
	TimeLeft    int64                       `json:"time_left"`
	Temperature MoveToNextStatusTemperature `json:"temperature,omitempty"`
}

type Status struct {
	Running     bool                     `json:"running"`
	Done        bool                     `json:"done"`
	PhaseNumber uint                     `json:"phase_number"`
	StartTime   time.Time                `json:"start_time"`
	EndTime     time.Time                `json:"end_time"`
	Next        MoveToNextStatus         `json:"next"`
	Heaters     []HeaterPhaseStatus      `json:"heaters"`
	Temperature []TemperaturePhaseStatus `json:"temperature"`
	GPIO        []GPIOPhaseStatus        `json:"gpio"`
	Errors      []string                 `json:"errors"`
}

type config struct {
	sensors map[string]*sensor
	heaters map[string]*heater
	outputs map[string]*output
	clock   Clock
	Config
}

func newConfig() *config {
	c := &config{
		sensors: map[string]*sensor{},
		heaters: map[string]*heater{},
		outputs: map[string]*output{},
		clock:   new(clock),
		Config:  Config{},
	}
	return c
}

func (c *config) UpdateSensors(sensors []Sensor) {
	for _, sensor := range sensors {
		c.sensors[sensor.ID()] = newSensor(sensor)
	}
	c.updateSensorsConfig()
}

func (c *config) UpdateHeaters(heaters []Heater) {
	for _, heater := range heaters {
		c.heaters[heater.ID()] = newHeater(heater)
	}
	c.updateHeaterConfig()
}

func (c *config) UpdateOutputs(outputs []Output) {
	for _, output := range outputs {
		c.outputs[output.ID()] = newOutput(output)
	}
	c.updateGPIOConfig()
}

func (c *config) SetPhaseNumber(number uint) {
	if number == c.PhaseNumber {
		return
	}

	c.Phases = resizeSlice(number, c.Phases)
	c.PhaseNumber = uint(len(c.Phases))

	// Update config
	c.updateGPIOConfig()
	c.updateSensorsConfig()
	c.updateHeaterConfig()
}

func (c *config) SetPhaseConfig(nb uint, conf PhaseConfig) error {
	if nb >= c.PhaseNumber {
		return ErrNoSuchPhase
	}

	if err := c.validateNextConfig(&conf.Next); err != nil {
		return err
	}

	if err := c.validateHeaterConfig(conf.Heaters); err != nil {
		return err
	}

	if err := c.validateGPIOConfig(conf.GPIO); err != nil {
		return err
	}

	c.Phases[nb] = conf
	return nil
}
func (c *config) SetGPIOGlobalConfig(conf []GPIOConfig) error {
	if err := c.validateGPIOConfig(conf); err != nil {
		return err
	}
	c.GlobalGPIO = conf
	return nil
}

func (c *config) GetConfig() Config {
	return c.Config
}

func (c *config) updateGPIOConfig() {
	size := len(c.outputs)
	// Update Global GPIO len
	c.GlobalGPIO = resizeSlice(size, c.GlobalGPIO)
	// Update len for each phase
	for i := range c.Phases {
		c.Phases[i].GPIO = resizeSlice(size, c.Phases[i].GPIO)
	}

	outputIDs := lo.Keys(c.outputs)
	// Check if there is any ID not associated with current outputs
	for i, elem := range c.GlobalGPIO {
		found := slices.Contains(outputIDs, elem.ID)
		// ID wasn't found
		if !found {
			// Clean it, later we will add appropriate id
			c.GlobalGPIO[i].ID = ""
		}
	}
	// And same for phases
	for i := range c.Phases {
		for j, elem := range c.Phases[i].GPIO {
			found := slices.Contains(outputIDs, elem.ID)
			// ID wasn't found
			if !found {
				// Clean it, later we will add appropriate id
				c.Phases[i].GPIO[j].ID = ""
			}
		}
	}

	for _, id := range outputIDs {
		// Do it for global config
		found := slices.ContainsFunc(c.GlobalGPIO, func(config GPIOConfig) bool {
			return config.ID == id
		})

		if !found {
			// ID not found
			// Place it on the first free slot
			idx := slices.IndexFunc(c.GlobalGPIO, func(config GPIOConfig) bool {
				return config.ID == ""
			})

			if idx != -1 {
				c.GlobalGPIO[idx].ID = id
			} else {
				// That shouldn't happen, however print it out
				log.Println("Free slot for OutputID not found")
			}
		}

		// And for phase configs
		for j := range c.Phases {
			// Do it for global config
			found := slices.ContainsFunc(c.Phases[j].GPIO, func(config GPIOConfig) bool {
				return config.ID == id
			})

			if !found {
				// ID not found
				// Place it on the first free slot
				idx := slices.IndexFunc(c.Phases[j].GPIO, func(config GPIOConfig) bool {
					return config.ID == ""
				})

				if idx != -1 {
					c.Phases[j].GPIO[idx].ID = id
				} else {
					// That shouldn't happen, however print it out
					log.Println("Free slot for OutputID not found")
				}
			}
		}
	}
}

func (c *config) updateHeaterConfig() {
	size := len(c.heaters)
	// Update slice for each Phase
	for i := range c.Phases {
		c.Phases[i].Heaters = resizeSlice(size, c.Phases[i].Heaters)
	}

	heaterIDs := lo.Keys(c.heaters)
	// Does each phase contains needed keys?
	for i := range c.Phases {
		for j, elem := range c.Phases[i].Heaters {
			found := slices.Contains(heaterIDs, elem.ID)
			// ID wasn't found
			if !found {
				// Clean it, later we will add appropriate id
				c.Phases[i].Heaters[j].ID = ""
			}
		}
	}

	// Make sure each ID can be found in Config
	for _, id := range heaterIDs {
		for j := range c.Phases {
			found := slices.ContainsFunc(c.Phases[j].Heaters, func(config HeaterPhaseConfig) bool {
				return config.ID == id
			})

			if !found {
				// ID not found
				// Place it on the first free slot
				idx := slices.IndexFunc(c.Phases[j].Heaters, func(config HeaterPhaseConfig) bool {
					return config.ID == ""
				})

				if idx != -1 {
					c.Phases[j].Heaters[idx].ID = id
				} else {
					// That shouldn't happen, however print it out
					log.Println("Free slot for HeaterID not found")
				}
			}
		}
	}
}

func (c *config) updateSensorsConfig() {
	size := len(c.sensors)
	sensorIds := lo.Keys(c.sensors)
	// Update slice for each Phase
	for i := range c.Phases {
		c.Phases[i].Next.Sensors = append(make([]string, 0, size), sensorIds...)
	}
}

func (c *config) validateHeaterConfig(heaters []HeaterPhaseConfig) error {
	if len(heaters) == 0 || (len(heaters) != len(c.heaters)) {
		return ErrHeaterConfigDiffersFromHeatersLen
	}

	ids := make(map[string]int)
	for _, heater := range heaters {
		if ids[heater.ID]++; ids[heater.ID] > 1 {
			return ErrDuplicatedID
		}
		if _, ok := c.heaters[heater.ID]; !ok {
			return ErrWrongHeaterID
		}
		if heater.Power > 100 || heater.Power < 0 {
			return ErrWrongHeaterPower
		}
	}
	return nil
}

func (c *config) validateGPIOConfig(conf []GPIOConfig) error {
	if len(conf) != len(c.outputs) {
		return ErrDifferentGPIOSConfig
	}

	ids := make(map[string]int)
	for _, con := range conf {
		// Check if there is not duplicated ID
		if ids[con.ID]++; ids[con.ID] > 1 {
			return ErrDuplicatedID
		}
		// Check if output is available
		if _, ok := c.outputs[con.ID]; !ok {
			return ErrWrongGpioID
		}
		// Check if sensor is available
		if _, ok := c.sensors[con.SensorID]; !ok {
			return ErrWrongSensorID
		}
	}
	return nil
}

func (c *config) validateNextConfig(next *MoveToNextConfig) error {
	switch next.Type {
	case ByTime:
		if next.TimeLeft <= 0 {
			return ErrByTimeWrongTime
		}
	case ByTemperature:
		if _, ok := c.sensors[next.SensorID]; !ok {
			return ErrByTemperatureWrongID
		}
	default:
		return ErrUnknownType
	}
	return nil
}

func (c *config) Validate() error {
	if len(c.sensors) == 0 {
		return ErrNoTSensors
	}

	if len(c.heaters) == 0 {
		return ErrNoHeaters
	}

	if err := c.validateGPIOConfig(c.Config.GlobalGPIO); err != nil {
		return err
	}

	for _, conf := range c.Config.Phases {
		if err := c.validateGPIOConfig(conf.GPIO); err != nil {
			return err
		}
		if err := c.validateHeaterConfig(conf.Heaters); err != nil {
			return err
		}
		if err := c.validateNextConfig(&conf.Next); err != nil {
			return err
		}
	}
	return nil
}

func resizeSlice[S constraints.Integer, T any](newSize S, s []T) []T {
	for i := S(len(s)); i < newSize; i++ {
		var t T
		s = append(s, t)
	}
	return s[:newSize]
}
