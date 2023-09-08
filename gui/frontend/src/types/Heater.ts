import { HeaterEnable } from '../../wailsjs/go/backend/Backend'
import { parameters } from '../../wailsjs/go/models'

export class Heater {
    heater: parameters.Heater;
    constructor(name: string, enabled: boolean) {
        this.heater = new parameters.Heater()
        this.heater.ID = name
        this.heater.enabled = enabled
    }

    get name() {
        return this.heater.ID
    }

    set enable(value: boolean) {
        this.heater.enabled = value
        HeaterEnable(this.heater.ID, this.heater.enabled)

    }

    get enable(): boolean {
        return this.heater.enabled
    }

}
