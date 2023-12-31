import { PTSetCorrection, PTSetSamples, PTEnable, PTSetName } from "../../wailsjs/go/backend/Backend";
import Parameter from "./Parameter";

export class PT100 {
    name: Parameter;
    id: string;
    correction: Parameter;
    samples: Parameter;
    temperature: string;
    private enable_: boolean;

    constructor(name: string, id: string, enabled:boolean, correction: number, samples: number, temperature: string = "") {
        this.id = id
        this.name = new Parameter(name, false, this.writeName)
        this.correction = new Parameter(correction, true, this.writeCorrection)
        this.samples = new Parameter(samples, false, this.writeSamples)
        this.temperature = temperature;
        this.enable_ = enabled
    }

    writeCorrection(value: number) {
        this.correction.value = value
        PTSetCorrection(this.id, value)
    }

    writeName(value: string) {
        this.name.value = value
        PTSetName(this.id, value)
    }

    writeSamples(value: number) {
        this.samples.value = value
        PTSetSamples(this.id, value)
    }

    set enable(value) {
        this.enable_ = value
        PTEnable(this.id, value)
    }

    get enable(): boolean {
        return this.enable_
    }
}
