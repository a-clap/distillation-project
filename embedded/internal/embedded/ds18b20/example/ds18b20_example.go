/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"fmt"
	"github.com/a-clap/iot/internal/embedded/ds18b20"
	"github.com/a-clap/iot/internal/embedded/logger"

	"go.uber.org/zap/zapcore"
	"time"
)

func main() {
	log := logger.NewDefaultLogger(zapcore.DebugLevel)
	ds, err := ds18b20.NewBus(ds18b20.WithOnewire())
	if err != nil {
		log.Panic(err)
	}

	ids, err := ds.IDs()
	if err != nil && len(ids) == 0 {
		log.Fatal(err)
	}

	sensor, _ := ds.NewSensor(ids[0])

	errs := sensor.Poll()
	if errs != nil {
		log.Fatal(err)
	}

	// Just to end this after time
	select {
	case <-time.After(10 * time.Second):
		_ = sensor.Close()
	}
	reads := sensor.GetReadings()

	for _, readings := range reads {
		fmt.Printf("id: %s, Temperature: %v. Time: %s, err: %v \n", readings.ID, readings.Temperature, readings.Stamp, readings.Error)
	}

	fmt.Println("finished")
}
