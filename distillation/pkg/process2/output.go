package process2

import (
	"fmt"
)

type Output interface {
	ID() string           // ID returns unique ID of interface
	Set(value bool) error // Set applies value to output
}

type output struct {
	Output  Output
	inRange bool
}

func newOutput(o Output) *output {
	return &output{
		Output:  o,
		inRange: false,
	}
}

func (o *output) ID() string {
	return o.Output.ID()
}

func (o *output) Set(v bool) error {
	err := o.Output.Set(v)
	if err != nil {
		// We don't allow any error from output
		return fmt.Errorf("%w: on ID: %v", ErrSetValue, o.ID())
	}
	return nil
}
