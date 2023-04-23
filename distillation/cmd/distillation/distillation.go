/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"log"
	"net"
	"time"
	
	"github.com/a-clap/distillation/pkg/distillation"
	"github.com/a-clap/embedded/pkg/embedded"
)

const EmbeddedAddr = "localhost:50001"

func main() {
	err := WaitForEmbedded(EmbeddedAddr, 30*time.Second)
	
	if err != nil {
		log.Fatalln("Couldn't connect to ", EmbeddedAddr)
	}
	
	heaterClient, err := embedded.NewHeaterRPCCLient(EmbeddedAddr, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer heaterClient.Close()
	
	ptClient, err := embedded.NewPTRPCClient(EmbeddedAddr, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer ptClient.Close()
	
	dsClient, err := embedded.NewDSRPCClient(EmbeddedAddr, 1*time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer dsClient.Close()
	
	gpioClient, err := embedded.NewGPIORPCClient(EmbeddedAddr, time.Second)
	if err != nil {
		log.Fatal(err)
	}
	defer gpioClient.Close()
	
	handler, err := distillation.NewRest(
		distillation.WithPT(ptClient),
		distillation.WithDS(dsClient),
		distillation.WithHeaters(heaterClient),
		distillation.WithGPIO(gpioClient),
		distillation.WithURL("localhost:8081"),
	)
	if err != nil {
		log.Fatalln(err)
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
			<-time.After(100 * time.Millisecond)
			continue
		}
		_ = conn.Close()
		break
	}
	
	return err
	
}
