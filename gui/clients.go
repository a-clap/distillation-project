package main

import (
	"fmt"
	"log"
	"time"

	"distillation/pkg/distillation"
	"gui/backend"
	"osservice"
)

func getopts(host string, distPort int, osPort int) []backend.Option {
	distAddr := fmt.Sprintf("%v:%v", host, distPort)

	heaterClient, err := distillation.NewHeaterRPCCLient(distAddr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	dsClient, err := distillation.NewDSRPCClient(distAddr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	ptClient, err := distillation.NewPTRPCClient(distAddr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	gpioClient, err := distillation.NewGPIORPCClient(distAddr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}
	phaseClient, err := distillation.NewProcessRPCClient(distAddr, time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	wifiClient, err := osservice.NewWifiClient(host, osPort, time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	return []backend.Option{
		backend.WithHeaterClient(heaterClient),
		backend.WithDSClient(dsClient),
		backend.WithPTClient(ptClient),
		backend.WithGPIOClient(gpioClient),
		backend.WithPhaseClient(phaseClient),
		backend.WithWifi(wifiClient),
	}
}
