package process2

import (
	"fmt"
)

type Heater interface {
	ID() string             // ID returns unique ID of interface
	SetPower(pwr int) error // SetPower set power (in %) for that Heater. 0% means Heater should be disabled
}

type heater struct {
	Heater Heater
}

func newHeater(h Heater) *heater {
	return &heater{
		Heater: h,
	}
}

func (h *heater) ID() string {
	return h.Heater.ID()
}

func (h *heater) SetPower(pwr int) error {
	err := h.Heater.SetPower(pwr)
	if err != nil {
		// We don't allow any error from heater
		return fmt.Errorf("%w: %v on ID: %v", err, ErrSetPower, h.ID())
	}
	return nil
}
