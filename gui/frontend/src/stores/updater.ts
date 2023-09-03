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

import {defineStore} from "pinia";
import {NotifyUpdate} from "../../wailsjs/go/backend/Events";
import {CheckUpdates, StartUpdate, StopUpdate} from "../../wailsjs/go/backend/Backend";
import {backend} from "../../wailsjs/go/models";

export const useUpdaterStore = defineStore('updater', {
    state: () => {
        return {
            updating: false as boolean,
            new_update: false as boolean,
            release: "" as string,
            state: 0 as number,
            downloading: 0 as number,
            installing: 0 as number,
            reboot: false as boolean,
            reboot_in: 0 as number,
        }
    },
    actions: {
        init() {
            NotifyUpdate().then((ev: string) => {
                return runtime.EventsOn(ev, (...args: any) => {
                    this.handle(...args);
                });
            })
        },
        checkUpdate(): Promise<void> {
            return new Promise((resolve, reject) => {
                CheckUpdates().then((data: backend.CheckUpdateData) => {
                    this.new_update = data.new_update
                    this.release = data.releases[0]
                    return data.error_code == 0 ? resolve() : reject(data.error_code)
                }).catch((error) => {
                    this.new_update = false
                    return reject(error)
                })
            })
        },
        startUpdate() {
            StartUpdate(this.release).catch((error) => {
                console.log(error)
            })
        },
        stopUpdate() {
            StopUpdate().catch((error) => {
                console.log(error)
            })
        },
        handle(...args: any) {
            let u = new backend.Update(args[0])
            this.updating = u.updating
            this.downloading = u.downloading
            this.installing = u.installing

            if (u.rebooting > 0 && !this.reboot) {
                this.reboot = true
                this.reboot_in = 0
                let reboot_timer = setInterval(() => {
                    this.reboot_in += 10
                    if (this.reboot_in == 100) {

                        clearInterval(reboot_timer)
                    }
                }, 1000)
            }

            console.log(u)
        }
    }
})
