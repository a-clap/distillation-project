/*
 * Copyright (c) 2023 a-clap. All rights reserved.
 * Use of this source code is governed by a MIT-style license that can be found in the LICENSE file.
 */

package distillation

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RoutesEnableHeater          = "/api/heater"
	RoutesGetAllHeaters         = "/api/heater"
	RoutesGetEnabledHeaters     = "/api/heater/enabled"
	RoutesConfigureHeater       = "/api/heater/enabled"
	RoutesGetDS                 = "/api/onewire"
	RoutesConfigureDS           = "/api/onewire"
	RoutesGetDSTemperatures     = "/api/onewire/temperatures"
	RoutesGetPT                 = "/api/pt100"
	RoutesConfigurePT           = "/api/pt100"
	RoutesGetPTTemperatures     = "/api/pt100/temperatures"
	RoutesGetGPIO               = "/api/gpio"
	RoutesConfigureGPIO         = "/api/gpio"
	RoutesProcessPhases         = "/api/phases"
	RoutesProcessConfigPhase    = "/api/phases/:id"
	RoutesProcessConfigValidate = "/api/process/validate"
	RoutesProcess               = "/api/process"
	RoutesProcessStatus         = "/api/process/status"
	RoutesProcessComponents     = "/api/process/components"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

type Rest struct {
	url string
	*gin.Engine
	*Distillation
}

func NewRest(url string, opts ...Option) (*Rest, error) {
	r := &Rest{
		url:    url,
		Engine: gin.Default(),
	}
	d, err := New(opts...)
	if err != nil {
		return nil, err
	}

	r.Distillation = d
	r.routes()

	return r, nil
}

func (r *Rest) Run() error {
	r.Distillation.Run()
	defer r.Distillation.Close()
	err := r.Engine.Run(r.url)

	return err
}

// routes configures default handlers for paths above
func (r *Rest) routes() {
	r.GET(RoutesGetAllHeaters, r.getAllHeaters())
	r.GET(RoutesGetEnabledHeaters, r.getEnabledHeaters())
	r.PUT(RoutesEnableHeater, r.enableHeater())
	r.PUT(RoutesConfigureHeater, r.configEnabledHeater())

	r.GET(RoutesGetDS, r.getDS())
	r.GET(RoutesGetDSTemperatures, r.getDSTemperatures())
	r.PUT(RoutesConfigureDS, r.configureDS())

	r.GET(RoutesGetPT, r.getPT())
	r.GET(RoutesGetPTTemperatures, r.getPTTemperatures())
	r.PUT(RoutesConfigurePT, r.configurePT())

	r.GET(RoutesGetGPIO, r.getGPIO())
	r.PUT(RoutesConfigureGPIO, r.configureGPIO())

	r.GET(RoutesProcessPhases, r.getPhaseCount())
	r.PUT(RoutesProcessPhases, r.configurePhaseCount())
	r.GET(RoutesProcessConfigPhase, r.getProcessConfig())
	r.PUT(RoutesProcessConfigPhase, r.setProcessConfig())

	r.GET(RoutesProcessConfigValidate, r.getConfigValidation())
	r.GET(RoutesProcessStatus, r.getProcessStatus())
	r.GET(RoutesProcessComponents, r.getComponents())
	r.PUT(RoutesProcess, r.configureProcess())

}

// common respond for whole rest API
func (*Rest) respond(ctx *gin.Context, code int, obj any) {
	ctx.JSON(code, obj)
}

func (r *Rest) configEnabledHeater() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.HeatersHandler == nil {
			e := &Error{
				Title:     "Failed to ConfigEnableHeater",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesConfigureHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		cfg := HeaterConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind HeaterConfig",
				Detail:    err.Error(),
				Instance:  RoutesConfigureHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		cfg, err := r.HeatersHandler.Configure(cfg)
		if err != nil {
			e := &Error{
				Title:     "Failed to Configure",
				Detail:    err.Error(),
				Instance:  RoutesConfigureHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, cfg)
	}
}
func (r *Rest) enableHeater() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.HeatersHandler == nil {
			e := &Error{
				Title:     "Failed to ConfigHeater",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesEnableHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		cfg := HeaterConfigGlobal{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind HeaterConfigGlobal",
				Detail:    err.Error(),
				Instance:  RoutesEnableHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}
		newCfg, err := r.HeatersHandler.ConfigureGlobal(cfg)
		if err != nil {
			e := &Error{
				Title:     "Failed to ConfigureGlobal",
				Detail:    err.Error(),
				Instance:  RoutesEnableHeater,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, newCfg)
	}
}

func (r *Rest) getAllHeaters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var heaters []HeaterConfigGlobal
		if r.HeatersHandler != nil {
			heaters = r.HeatersHandler.ConfigsGlobal()
		}
		if len(heaters) == 0 {
			e := &Error{
				Title:     "Failed to ConfigsGlobal",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetAllHeaters,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, heaters)
	}
}

func (r *Rest) getEnabledHeaters() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.HeatersHandler == nil {
			e := &Error{
				Title:     "Failed to GetEnabledHeaters",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetEnabledHeaters,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		heaters := r.HeatersHandler.Configs()
		r.respond(ctx, http.StatusOK, heaters)
	}
}

func (r *Rest) getDS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var sensors []DSConfig
		if r.Distillation.DSHandler != nil {
			sensors = r.Distillation.DSHandler.GetSensors()
		}
		if len(sensors) == 0 {
			e := &Error{
				Title:     "Failed to GetSensors",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetDS,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, sensors)
	}
}

func (r *Rest) getDSTemperatures() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var temperatures []DSTemperature
		if r.Distillation.DSHandler != nil {
			temperatures = r.Distillation.DSHandler.Temperatures()
		}
		if len(temperatures) == 0 {
			e := &Error{
				Title:     "Failed to get Temperatures",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetDSTemperatures,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, temperatures)
	}
}

func (r *Rest) configureDS() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.Distillation.DSHandler == nil {
			e := &Error{
				Title:     "Failed to ConfigureSensor",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesConfigureDS,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		cfg := DSConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind DSConfig",
				Detail:    err.Error(),
				Instance:  RoutesConfigureDS,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		newcfg, err := r.Distillation.DSHandler.ConfigureSensor(cfg)
		if err != nil {
			e := &Error{
				Title:     "Failed to ConfigureSensor",
				Detail:    err.Error(),
				Instance:  RoutesConfigureDS,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, newcfg)
	}
}

func (r *Rest) getPT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var sensors []PTConfig
		if r.Distillation.PTHandler != nil {
			sensors = r.Distillation.PTHandler.GetSensors()
		}
		if len(sensors) == 0 {
			e := &Error{
				Title:     "Failed to GetSensors",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetPT,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, sensors)
	}
}

func (r *Rest) getPTTemperatures() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var temperatures []PTTemperature
		if r.Distillation.PTHandler != nil {
			temperatures = r.Distillation.PTHandler.Temperatures()
		}
		if len(temperatures) == 0 {
			e := &Error{
				Title:     "Failed to get Temperatures",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetPTTemperatures,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, temperatures)
	}
}

func (r *Rest) configurePT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.Distillation.PTHandler == nil {
			e := &Error{
				Title:     "Failed to ConfigurePT",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesConfigurePT,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}

		cfg := PTConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind PTConfig",
				Detail:    err.Error(),
				Instance:  RoutesConfigurePT,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		c, err := r.Distillation.PTHandler.Configure(cfg)
		if err != nil {
			e := &Error{
				Title:     "Failed to Configure",
				Detail:    err.Error(),
				Instance:  RoutesConfigurePT,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, c)
	}
}

func (r *Rest) getGPIO() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var gpios []GPIOConfig
		if r.Distillation.GPIOHandler != nil {
			gpios = r.Distillation.GPIOHandler.Config()
		}
		if len(gpios) == 0 {
			e := &Error{
				Title:     "Failed to get Config",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesGetGPIO,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, gpios)
	}
}
func (r *Rest) configureGPIO() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if r.Distillation.GPIOHandler == nil {
			e := &Error{
				Title:     "Failed to ConfigGPIO",
				Detail:    ErrNotImplemented.Error(),
				Instance:  RoutesConfigureGPIO,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}

		cfg := GPIOConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind GPIOConfig",
				Detail:    err.Error(),
				Instance:  RoutesConfigureGPIO,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		newCfg, err := r.Distillation.GPIOHandler.Configure(cfg)
		if err != nil {
			e := &Error{
				Title:     "Failed to Configure",
				Detail:    err.Error(),
				Instance:  RoutesConfigureGPIO,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, newCfg)
	}
}

func (r *Rest) configureProcess() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := ProcessConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind ProcessConfig",
				Detail:    err.Error(),
				Instance:  RoutesProcess,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}
		if err := r.Distillation.ConfigureProcess(cfg); err != nil {
			e := &Error{
				Title:     "Failed ConfigureProcess",
				Detail:    err.Error(),
				Instance:  RoutesProcess,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusInternalServerError, e)
			return
		}
		r.respond(ctx, http.StatusOK, cfg)
	}
}

func (r *Rest) getProcessStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r.respond(ctx, http.StatusOK, r.Distillation.Status())
	}
}

func (r *Rest) getComponents() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r.respond(ctx, http.StatusOK, r.Distillation.Process.Components())
	}
}

func (r *Rest) getConfigValidation() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		v := r.Distillation.ValidateConfig()
		r.respond(ctx, http.StatusOK, v)
	}
}

func (r *Rest) getPhaseCount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := r.Distillation.Process.GetConfig()
		s := ProcessPhaseCount{PhaseNumber: cfg.PhaseNumber}
		r.respond(ctx, http.StatusOK, s)
	}
}

func (r *Rest) configurePhaseCount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cfg := ProcessPhaseCount{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind ProcessPhaseCount",
				Detail:    err.Error(),
				Instance:  RoutesProcessPhases,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		if err := r.Distillation.Process.SetPhases(cfg.PhaseNumber); err != nil {
			e := &Error{
				Title:     "Failed to SetPhases",
				Detail:    err.Error(),
				Instance:  RoutesProcessPhases,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		config := r.Distillation.Process.GetConfig()
		s := ProcessPhaseCount{PhaseNumber: config.PhaseNumber}
		r.respond(ctx, http.StatusOK, s)

	}
}

func (r *Rest) getProcessConfig() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			e := &Error{
				Title:     "Failed to parse \"id\"",
				Detail:    err.Error(),
				Instance:  RoutesProcessConfigPhase,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}
		cfg := r.Distillation.Process.GetConfig()
		if int(id) >= cfg.PhaseNumber {
			e := &Error{
				Title:     "Phase doesn't exist",
				Detail:    fmt.Errorf("requested phase %v doesn't exist", id).Error(),
				Instance:  RoutesProcessConfigPhase,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}
		r.respond(ctx, http.StatusOK, cfg.Phases[id])
	}
}

func (r *Rest) setProcessConfig() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		param := ctx.Param("id")
		id, err := strconv.ParseInt(param, 10, 32)
		if err != nil {
			e := &Error{
				Title:     "Failed to parse \"id\"",
				Detail:    err.Error(),
				Instance:  RoutesProcessConfigPhase,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}
		cfg := ProcessPhaseConfig{}
		if err := ctx.ShouldBind(&cfg); err != nil {
			e := &Error{
				Title:     "Failed to bind ProcessPhaseConfig",
				Detail:    err.Error(),
				Instance:  RoutesProcessConfigPhase,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		if err := r.Distillation.configurePhase(int(id), cfg); err != nil {
			e := &Error{
				Title:     "Failed to configurePhase",
				Detail:    err.Error(),
				Instance:  RoutesProcessConfigPhase,
				Timestamp: time.Now(),
			}
			r.respond(ctx, http.StatusBadRequest, e)
			return
		}

		config := r.Distillation.Process.GetConfig()
		r.respond(ctx, http.StatusOK, config.Phases[id])
	}
}
