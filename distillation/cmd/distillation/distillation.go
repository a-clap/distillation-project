/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"flag"
	"log"
	"net"
	"strconv"
	"time"

	"distillation/pkg/distillation"
	"embedded/pkg/embedded"
)

var (
	embeddedPort = flag.Int("embedded_port", 50001, "embedded server port")
	embeddedRest = flag.Bool("embedded_rest", false, "embedded uses REST API")
	port         = flag.Int("port", 50002, "the server port")
	rest         = flag.Bool("rest", false, "use REST API instead of gRPC")
)

type embeddedHeaterClient interface {
	Get() ([]embedded.HeaterConfig, error)
	Configure(heater embedded.HeaterConfig) (embedded.HeaterConfig, error)
}

type embeddedPTClient interface {
	Get() ([]embedded.PTSensorConfig, error)
	Configure(s embedded.PTSensorConfig) (embedded.PTSensorConfig, error)
	Temperatures() ([]embedded.PTTemperature, error)
}

type embeddedGPIOClient interface {
	Get() ([]embedded.GPIOConfig, error)
	Configure(c embedded.GPIOConfig) (embedded.GPIOConfig, error)
}

type embeddedDSClient interface {
	Get() ([]embedded.DSSensorConfig, error)
	Configure(s embedded.DSSensorConfig) (embedded.DSSensorConfig, error)
	Temperatures() ([]embedded.DSTemperature, error)
}

type handler interface {
	Run() error
	Close()
}

func main() {
	flag.Parse()
	// Embedded clients
	var heaterClient embeddedHeaterClient
	var gpioClient embeddedGPIOClient
	var ptClient embeddedPTClient
	var dsClient embeddedDSClient

	const timeout = time.Second
	embeddedAddr := ":" + strconv.FormatInt(int64(*embeddedPort), 10)
	distillationAddr := ":" + strconv.FormatInt(int64(*port), 10)

	if err := WaitForEmbedded(embeddedAddr, 30*time.Second); err != nil {
		log.Fatalf("Couldn't connect to %v, err as follows: %v", embeddedAddr, err)
	}

	// Embedded uses rest
	if *embeddedRest {
		log.Println("Using REST clients to communicate with embedded...")
		embeddedAddr = "http://" + embeddedAddr
		hClient := embedded.NewHeaterClient(embeddedAddr, timeout)
		gClient := embedded.NewGPIOClient(embeddedAddr, timeout)
		pClient := embedded.NewPTClient(embeddedAddr, timeout)
		dClient := embedded.NewDS18B20Client(embeddedAddr, timeout)

		ptClient = pClient
		gpioClient = gClient
		heaterClient = hClient
		dsClient = dClient
	} else {
		log.Println("Using gRPC clients to communicate with embedded...")
		hClient, err := embedded.NewHeaterRPCCLient(embeddedAddr, 1*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		defer hClient.Close()

		pClient, err := embedded.NewPTRPCClient(embeddedAddr, 1*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		defer pClient.Close()

		dClient, err := embedded.NewDSRPCClient(embeddedAddr, 1*time.Second)
		if err != nil {
			log.Fatal(err)
		}
		defer dClient.Close()

		gClient, err := embedded.NewGPIORPCClient(embeddedAddr, time.Second)
		if err != nil {
			log.Fatal(err)
		}
		defer gClient.Close()

		ptClient = pClient
		gpioClient = gClient
		heaterClient = hClient
		dsClient = dClient
	}

	opts := []distillation.Option{
		distillation.WithPT(ptClient),
		distillation.WithDS(dsClient),
		distillation.WithHeaters(heaterClient),
		distillation.WithGPIO(gpioClient),
	}

	var err error
	var handler handler
	if *rest {
		log.Println("Running distillation as REST server on ", distillationAddr)
		handler, err = distillation.NewRest(distillationAddr, opts...)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		log.Println("Running distillation as RPC server on ", distillationAddr)
		handler, err = distillation.NewRPC(distillationAddr, opts...)
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = handler.Run()
	log.Println(err)
}

func WaitForEmbedded(addr string, timeout time.Duration) error {
	var err error
	deadLine := time.Now().Add(timeout)
	for deadLine.After(time.Now()) {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			log.Println("failed ", err)
			<-time.After(100 * time.Millisecond)
			continue
		}
		_ = conn.Close()
		break
	}

	return err

}
