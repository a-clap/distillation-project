/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/a-clap/iot/internal/distillation"
	"github.com/a-clap/iot/internal/embedded"
)

const EmbeddedAddr = "http://localhost:8080"

func main() {
	distil, err := distillation.New(distillation.WithPT(&PTHandler{}))
	if err != nil {
		log.Fatalln(err)
	}
	pts := distil.PTHandler.GetSensors()
	for _, pt := range pts {
		pt.Enabled = true
		pt.Samples = 10
		if _, err := distil.PTHandler.Configure(pt); err != nil {
			log.Println(err)
		}
	}

	err = distil.Run("localhost:8081")
	log.Println(err)
}

type PTHandler struct {
}

func (p *PTHandler) Get() ([]embedded.PTSensorConfig, error) {
	ctx := context.Background()
	timeout := 1 * time.Second
	reqContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ret := []embedded.PTSensorConfig{}
	r, err := http.NewRequestWithContext(reqContext, http.MethodGet, EmbeddedAddr+embedded.RoutesGetPT100Sensors, nil)
	if err != nil {
		return ret, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return ret, err
	}

	err = json.NewDecoder(res.Body).Decode(&ret)
	return ret, err
}

func (p *PTHandler) Configure(cfg embedded.PTSensorConfig) (embedded.PTSensorConfig, error) {
	ctx := context.Background()
	timeout := 1 * time.Second
	reqContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	c := embedded.PTSensorConfig{}
	b, err := json.Marshal(&cfg)
	if err != nil {
		return c, err
	}
	byteReader := bytes.NewReader(b)
	r, err := http.NewRequestWithContext(reqContext, http.MethodPut, EmbeddedAddr+embedded.RoutesConfigPT100Sensor, byteReader)
	if err != nil {
		return c, err
	}
	r.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return c, err
	}

	if res.StatusCode == 200 {
		return c, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(body, &c); err != nil {
		return c, fmt.Errorf("%w: %s", err, errors.New(string(body)))
	}

	return c, nil
}

func (p *PTHandler) Temperatures() ([]embedded.PTTemperature, error) {
	ctx := context.Background()
	timeout := 1 * time.Second
	reqContext, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ret := []embedded.PTTemperature{}
	r, err := http.NewRequestWithContext(reqContext, http.MethodGet, EmbeddedAddr+embedded.RoutesGetPT100Temperatures, nil)
	if err != nil {
		return ret, err
	}
	res, err := http.DefaultClient.Do(r)
	if err != nil {
		return ret, err
	}

	err = json.NewDecoder(res.Body).Decode(&ret)
	return ret, err
}
