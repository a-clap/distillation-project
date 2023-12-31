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

import { defineStore } from "pinia";
import { PhasesDisable, PhasesEnable, PhasesMoveToNext, PhasesValidateConfig } from "../../wailsjs/go/backend/Backend";
import { ProcessListener } from "../types/ProcessListener";
import { backend, distillation } from "../../wailsjs/go/models";
import { useNameStore } from "./names";
import { ErrorListener } from "../types/ErrorListener";
import { AppErrorCodes } from "./error_codes";

export interface Button {
    is_enabled: boolean;
    enable: Function;
}

export interface Heater {
    id: string;
    pwr: string;
}

export interface Sensor {
    id: string;
    temperature: string;
}

export interface Output {
    id: string;
    state: boolean;
}


function padTo2Digits(num: number) {
    return num.toString().padStart(2, '0');
}

function formatDate(date: Date) {
    return (
        [
            padTo2Digits(date.getHours()),
            padTo2Digits(date.getMinutes()),
            padTo2Digits(date.getSeconds()),
        ].join(':')
    );
}

export const useProcessStore = defineStore('process', {
    state: () => {
        return {
            enable: {} as Button,
            moveToNext: {} as Button,
            disable: {} as Button,
            heaters: [] as Heater[],
            sensors: [] as Sensor[],
            outputs: [] as Output[],
            is_valid: false,
            running: false,
            show_status: false,
            start_time: "",
            end_time: "",
            current_phase: "",
            current_type_time: false,
            phase_timeleft: "",
            phase_sensor: "",
            phase_sensor_threshold: "",
            names: useNameStore(),
        }
    },
    actions: {
        init() {
            this.enable.enable = this.processEnable
            this.moveToNext.enable = this.processMoveToNext
            this.disable.enable = this.processDisable

            this.updateButtons()
            ProcessListener.subscribeValidate(this.onValidate)
            ProcessListener.subscribeStatus(this.onStatus)
        },

        reload() {
            PhasesValidateConfig()
        },

        onStatus(v: backend.ProcessStatus) {
            this.running = v.running
            if (v.running || v.done) {
                // Current phase from backend starts from 0
                this.current_phase = (v.phase_number + 1).toString()
                // Time
                this.current_type_time = v.next.type == 0
                this.phase_timeleft = v.next.time_left.toString()
                // type == 1 is temperature
                if (v.next.type == 1 && v.next.temperature) {
                    this.phase_sensor_threshold = v.next.temperature?.sensor_threshold.toFixed(2).toString()
                    let id = v.next.temperature?.sensor_id.toString()
                    let [name, got] = this.names.id_to_name(id)
                    if (got) {
                        id = name
                    } else {
                        ErrorListener.sendError(AppErrorCodes.SensorIDNotFound)
                    }

                    this.phase_sensor = id
                }
                this.start_time = formatDate(new Date(v.unix_start_time * 1000))
                if (v.done) {
                    this.end_time = formatDate(new Date(v.unix_end_time * 1000))
                } else {
                    this.end_time = ""
                }

                let heaters: Heater[] = []
                v.heaters.forEach((v) => {
                    let h: Heater = { id: v.ID, pwr: v.power.toString() }
                    heaters.push(h)
                })

                this.heaters = heaters.sort((a: Heater, b: Heater) => {
                    if (a.id > b.id) {
                        return 1
                    }
                    if (a.id < b.id) {
                        return -1
                    }
                    return 0
                })

                let outputs: Output[] = []
                v.gpio.forEach((v) => {
                    let o: Output = { id: v.id, state: v.state }
                    outputs.push(o)
                })

                this.outputs = outputs.sort((a: Output, b: Output) => {
                    if (a.id > b.id) {
                        return 1
                    }
                    if (a.id < b.id) {
                        return -1
                    }
                    return 0
                })

                let sensors: Sensor[] = []
                v.temperature.forEach((v) => {
                    let [name, got] = this.names.id_to_name(v.ID)
                    if (got) {
                        v.ID = name
                    } else {
                        ErrorListener.sendError(AppErrorCodes.SensorIDNotFound)
                    }
                    let s: Sensor = { id: v.ID, temperature: v.temperature.toFixed(2).toString() }
                    sensors.push(s)
                })

                this.sensors = sensors.sort((a: Sensor, b: Sensor) => {
                    if (a.id > b.id) {
                        return 1
                    }
                    if (a.id < b.id) {
                        return -1
                    }
                    return 0
                })


            }
            this.updateButtons()
            this.show_status = v.running || v.done
        },

        updateButtons() {
            this.moveToNext.is_enabled = this.running
            this.disable.is_enabled = this.running

            this.enable.is_enabled = !this.running && this.is_valid
        },

        onValidate(v: distillation.ProcessConfigValidation) {
            this.is_valid = v.valid
            this.updateButtons()
        },

        processEnable() {
            PhasesEnable()
        },
        processMoveToNext() {
            PhasesMoveToNext()
        },
        processDisable() {
            PhasesDisable()
        },
    }
})
