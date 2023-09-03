import { defineStore } from "pinia";
import {NotifyDSConfig, NotifyUpdate} from "../../wailsjs/go/backend/Events";
import {parameters} from "../../wailsjs/go/models";

export const useUpdaterStore = defineStore('updater', {
    state: () => {
        return {
            updating: true as boolean,
            newUpdate: false as boolean,
            state: 0 as number,
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

        checkUpdate() {
            setTimeout( () => { this.updating = !this.updating }, 1000 );
        },

        handle(...args: any){
            console.log(args[0])
        }
    }
})
