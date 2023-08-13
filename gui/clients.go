package main

import (
	"fmt"
	"log"
	"time"

	"distillation/pkg/distillation"
	"gui/backend"
	"osservice"
)

const (
	defaultTimeout = time.Second
)

func getopts(host string, distPort int, osPort int) []backend.Option {
	distAddr := fmt.Sprintf("%v:%v", host, distPort)

	heaterClient, err := distillation.NewHeaterRPCCLient(distAddr, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}
	dsClient, err := distillation.NewDSRPCClient(distAddr, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	ptClient, err := distillation.NewPTRPCClient(distAddr, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}
	gpioClient, err := distillation.NewGPIORPCClient(distAddr, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}
	phaseClient, err := distillation.NewProcessRPCClient(distAddr, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	wifiClient, err := osservice.NewWifiClient(host, osPort, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	// Time operations, especially NTP takes a lot of time
	timeClient, err := osservice.NewTimeClient(host, osPort, 10*time.Second)
	if err != nil {
		log.Fatalln(err)
	}

	// Time operations, especially NTP takes a lot of time
	netClient, err := osservice.NewNetClient(host, osPort, defaultTimeout)
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
		backend.WithTime(timeClient),
		backend.WithNet(netClient),
	}
}
