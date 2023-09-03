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
	"fmt"
	"log"
	"osservice"
	"time"

	"gui/backend"
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

	saverClient, err := osservice.NewStoreClient(host, osPort, defaultTimeout)
	if err != nil {
		log.Fatalln(err)
	}

	updateClient, err := osservice.NewUpdateClient(host, osPort, defaultTimeout)
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
		backend.WithLoadSaver(saverClient),
		backend.WithUpdate(updateClient),
	}
}
