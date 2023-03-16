/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/a-clap/iot/internal/distillation"
	"github.com/a-clap/iot/internal/embedded"
)

const EmbeddedAddr = "http://localhost:8080"

func main() {
	err := WaitForEmbedded(EmbeddedAddr, 30*time.Second)
	if err != nil {
		log.Fatalln("Couldn't connect to ", EmbeddedAddr)
	}

	handler, err := distillation.New(
		distillation.WithPT(embedded.NewPTClient(EmbeddedAddr, 1*time.Second)),
		distillation.WithDS(embedded.NewDS18B20Client(EmbeddedAddr, 1*time.Second)),
		distillation.WithHeaters(embedded.NewHeaterClient(EmbeddedAddr, 1*time.Second)),
		distillation.WithGPIO(embedded.NewGPIOClient(EmbeddedAddr, 1*time.Second)),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = handler.Run("localhost:8081")
	log.Println(err)
}

func WaitForEmbedded(addr string, timeout time.Duration) error {
	var err error
	deadLine := time.Now().Add(timeout)
	for deadLine.After(time.Now()) {
		_, err := http.Get(addr)
		if err == nil {
			break
		}
		<-time.After(100 * time.Millisecond)
	}

	return err

}
