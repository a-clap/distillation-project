// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {distillation} from '../models';
import {gpio} from '../models';
import {process} from '../models';
import {backend} from '../models';
import {parameters} from '../models';

export function DistillationProcessStatus():Promise<distillation.ProcessStatus>;

export function GPIOActiveLevel():Promise<gpio.ActiveLevel>;

export function GPIOPhaseConfig():Promise<process.GPIOConfig>;

export function GPIOPhaseStatus():Promise<process.GPIOPhaseStatus>;

export function HeaterPhaseConfig():Promise<process.HeaterPhaseConfig>;

export function HeaterPhaseStatus():Promise<process.HeaterPhaseStatus>;

export function MoveToNextConfig():Promise<process.MoveToNextConfig>;

export function MoveToNextStatus():Promise<process.MoveToNextStatus>;

export function MoveToNextStatusTemperature():Promise<process.MoveToNextStatusTemperature>;

export function ProcessConfigValidation():Promise<distillation.ProcessConfigValidation>;

export function ProcessStatus():Promise<backend.ProcessStatus>;

export function Temperature():Promise<parameters.Temperature>;

export function TemperatureErrorCodeEmptyBuffer():Promise<number>;

export function TemperatureErrorCodeInternal():Promise<number>;

export function TemperatureErrorCodeWrongID():Promise<number>;

export function TemperaturePhaseStatus():Promise<process.TemperaturePhaseStatus>;

export function Update():Promise<backend.Update>;
