package main

import (
	"embedded/pkg/embedded"
	"embedded/pkg/embeddedmock"
	"embedded/pkg/gpio"
)

func setupLogging() {
}

func getOpts(configPath string) ([]embedded.Option, []error) {
	ptIds := []string{"PT_1", "PT_2", "PT_3"}
	pts := make([]embedded.PTSensor, len(ptIds))
	for i, id := range ptIds {
		pts[i] = embeddedmock.NewPT(id)
	}

	dsIds := []struct {
		bus, id string
	}{
		{
			bus: "1",
			id:  "ds_1",
		},
		{
			bus: "1",
			id:  "ds_2",
		},
	}
	dss := make([]embedded.DSSensor, len(dsIds))
	for i, id := range dsIds {
		dss[i] = embeddedmock.NewDS(id.bus, id.id)
	}

	heaterIds := []string{"heater_1", "heater_2", "heater_3"}
	heaters := make(map[string]embedded.Heater, len(heaterIds))
	for _, id := range heaterIds {
		heaters[id] = embeddedmock.NewHeater()
	}

	gpioIds := []struct {
		id    string
		state bool
		dir   gpio.Direction
	}{
		{
			id:    "gpio_1",
			state: false,
			dir:   gpio.DirInput,
		}, {
			id:    "gpio_2",
			state: true,
			dir:   gpio.DirOutput,
		},
	}
	gpios := make([]embedded.GPIO, len(gpioIds))
	for i, id := range gpioIds {
		gpios[i] = embeddedmock.NewGPIO(id.id, id.state, id.dir)
	}

	return []embedded.Option{
		embedded.WithPT(pts),
		embedded.WithDS18B20(dss),
		embedded.WithHeaters(heaters),
		embedded.WithGPIOs(gpios),
	}, nil
}
