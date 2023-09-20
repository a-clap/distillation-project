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
import {NotifyUpdate, NotifyUpdateNextState, NotifyUpdateStatus} from "../../wailsjs/go/backend/Events";
import {CheckUpdates, ContinueUpdate, MoveToNextState, StartUpdate, StopUpdate} from "../../wailsjs/go/backend/Backend";
import {backend} from "../../wailsjs/go/models";
import {ErrorListener} from "../types/ErrorListener";
import {AppErrorCodes} from "./error_codes";

enum MenderStatus {
    Downloading = 1,
    PauseBeforeInstalling,
    Installing,
    PauseBeforeRebooting,
    Rebooting,
    PauseBeforeCommitting,
    Success,
    Failure,
    AlreadyInstalled,
}


export const useUpdaterStore = defineStore('updater', {
    state: () => {
        return {
            updating: false as boolean,
            new_update: false as boolean,
            releases: [] as string[],
            downloading: 0 as number,
            installing: 0 as number,
            reboot: false as boolean,
            rebooting: 0 as number,
            commit: false as boolean,
            message: "" as string,
        }
    },
    actions: {
        init() {
            NotifyUpdate().then((ev: string) => {
                return runtime.EventsOn(ev, (...args: any) => {
                    this.update(...args);
                })
            })
            NotifyUpdateStatus().then((ev: string) => {
                return runtime.EventsOn(ev, (...args: any) => {
                    this.updateStatus(...args);
                })
            })
            NotifyUpdateNextState().then((ev: string) => {
                return runtime.EventsOn(ev, (...args: any) => {
                    this.nextState(...args);
                })
            })

            ContinueUpdate().catch((error) => {
                console.log(error)
            })
        },
        checkUpdate(): Promise<void> {
            return new Promise((resolve, reject) => {
                CheckUpdates().then((data: backend.UpdateData) => {
                    this.releases = data.releases
                    this.new_update = this.releases.length > 0
                    return data.error_code == 0 ? resolve() : reject(data.error_code)
                }).catch((error) => {
                    this.new_update = false
                    return reject(error)
                })
            })
        },
        startUpdate() {
            StartUpdate(this.releases[0]).catch((error) => {
                console.log(error)
            })
        },
        stopUpdate() {
            StopUpdate().catch((error): void => {
                console.log(error)
            })
        },

        updateStatus(...args: any): void {
            let status = new backend.UpdateStateStatus(args[0])
            if (status.state == MenderStatus.Downloading) {
                this.downloading = status.progress
            } else if (status.state == MenderStatus.Installing) {
                this.downloading = 100
                this.installing = status.progress
            }

        },
        nextState(...args: any): void {
            let next = new backend.UpdateNextState(args[0])
            if (next.state == MenderStatus.Rebooting) {
                if (!this.reboot) {
                    this.reboot = true
                    this.rebooting = 0
                    let reboot_timer = setInterval(() => {
                        this.rebooting += 10
                        if (this.rebooting == 100) {
                            clearInterval(reboot_timer)
                            MoveToNextState(true).catch((error): void => {
                                console.log(error)
                            })
                        }
                    }, 1000)
                }
            } else if (next.state == MenderStatus.Success) {
                this.message = "blabla"
                this.commit = true
            } else if (next.state == MenderStatus.Failure || next.state == MenderStatus.AlreadyInstalled) {
                let err = next.state == MenderStatus.Failure ? AppErrorCodes.UpdateFail : AppErrorCodes.UpdateAlreadyInstalled;
                ErrorListener.sendError(err)
            } else {
                MoveToNextState(true).catch((error) => {
                    console.log(error)
                })
            }
        },
        cleanup(): void {
            this.releases = []
            this.commit = false
            this.updating = false
        },
        submit(): void {
            MoveToNextState(true).catch((error): void => {
                console.log(error)
            })
            this.cleanup()
        },
        update(...args: any): void {
            let u = new backend.Update(args[0])
            this.updating = u.updating

            if (!this.updating) {
                this.cleanup()
            }
        }
    }
})
