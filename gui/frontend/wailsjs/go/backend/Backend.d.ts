// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {parameters} from '../models';
import {backend} from '../models';
import {distillation} from '../models';
import {process} from '../models';

export function DSEnable(arg1:string,arg2:boolean):Promise<void>;

export function DSGet():Promise<Array<parameters.DS>>;

export function DSSetCorrection(arg1:string,arg2:number):Promise<void>;

export function DSSetName(arg1:string,arg2:string):Promise<void>;

export function DSSetResolution(arg1:string,arg2:number):Promise<void>;

export function DSSetSamples(arg1:string,arg2:number):Promise<void>;

export function Events():Promise<backend.Events>;

export function GPIOGet():Promise<Array<parameters.GPIO>>;

export function GPIOSetActiveLevel(arg1:string,arg2:number):Promise<void>;

export function GPIOSetState(arg1:string,arg2:boolean):Promise<void>;

export function HeaterEnable(arg1:string,arg2:boolean):Promise<void>;

export function HeatersGet():Promise<Array<parameters.Heater>>;

export function LoadParameters():Promise<void>;

export function PTEnable(arg1:string,arg2:boolean):Promise<void>;

export function PTGet():Promise<Array<parameters.PT>>;

export function PTSetCorrection(arg1:string,arg2:number):Promise<void>;

export function PTSetName(arg1:string,arg2:string):Promise<void>;

export function PTSetSamples(arg1:string,arg2:number):Promise<void>;

export function PhasesDisable():Promise<void>;

export function PhasesEnable():Promise<void>;

export function PhasesGetGlobalConfig():Promise<any>;

export function PhasesGetPhaseConfigs():Promise<Array<distillation.ProcessPhaseConfig>>;

export function PhasesGetPhaseCount():Promise<any>;

export function PhasesMoveToNext():Promise<void>;

export function PhasesSetConfig(arg1:number,arg2:distillation.ProcessPhaseConfig):Promise<void>;

export function PhasesSetGlobalGPIO(arg1:Array<process.GPIOConfig>):Promise<void>;

export function PhasesSetPhaseCount(arg1:number):Promise<void>;

export function PhasesValidateConfig():Promise<void>;

export function SaveParameters():Promise<void>;

export function WifiAPList():Promise<Array<string>>;
