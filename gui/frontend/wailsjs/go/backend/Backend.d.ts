// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {backend} from '../models';
import {parameters} from '../models';
import {mender} from '../models';
import {process} from '../models';
import {distillation} from '../models';

export function CheckUpdates():Promise<backend.CheckUpdateData>;

export function Commit(arg1:boolean):Promise<void>;

export function DSEnable(arg1:string,arg2:boolean):Promise<void>;

export function DSGet():Promise<Array<parameters.DS>>;

export function DSSetCorrection(arg1:string,arg2:number):Promise<void>;

export function DSSetName(arg1:string,arg2:string):Promise<void>;

export function DSSetResolution(arg1:string,arg2:number):Promise<void>;

export function DSSetSamples(arg1:string,arg2:number):Promise<void>;

export function Error(arg1:Error):Promise<void>;

export function Events():Promise<backend.Events>;

export function GPIOGet():Promise<Array<parameters.GPIO>>;

export function GPIOSetActiveLevel(arg1:string,arg2:number):Promise<void>;

export function GPIOSetState(arg1:string,arg2:boolean):Promise<void>;

export function HeaterEnable(arg1:string,arg2:boolean):Promise<void>;

export function HeatersGet():Promise<Array<parameters.Heater>>;

export function ListInterfaces():Promise<Array<backend.NetInterface>>;

export function LoadParameters():Promise<void>;

export function NTPGet():Promise<boolean>;

export function NTPSet(arg1:boolean):Promise<void>;

export function NextState(arg1:mender.DeploymentStatus):Promise<boolean>;

export function Now():Promise<number>;

export function PTEnable(arg1:string,arg2:boolean):Promise<void>;

export function PTGet():Promise<Array<parameters.PT>>;

export function PTSetCorrection(arg1:string,arg2:number):Promise<void>;

export function PTSetName(arg1:string,arg2:string):Promise<void>;

export function PTSetSamples(arg1:string,arg2:number):Promise<void>;

export function PhasesDisable():Promise<void>;

export function PhasesEnable():Promise<void>;

export function PhasesGetGlobalConfig():Promise<process.Config>;

export function PhasesGetPhaseConfigs():Promise<Array<distillation.ProcessPhaseConfig>>;

export function PhasesGetPhaseCount():Promise<distillation.ProcessPhaseCount>;

export function PhasesMoveToNext():Promise<void>;

export function PhasesSetConfig(arg1:number,arg2:distillation.ProcessPhaseConfig):Promise<void>;

export function PhasesSetGlobalGPIO(arg1:Array<process.GPIOConfig>):Promise<void>;

export function PhasesSetPhaseCount(arg1:number):Promise<void>;

export function PhasesValidateConfig():Promise<void>;

export function Reboot(arg1:boolean):Promise<void>;

export function SaveParameters():Promise<void>;

export function StartUpdate(arg1:string):Promise<void>;

export function StopUpdate():Promise<void>;

export function TimeSet(arg1:number):Promise<void>;

export function Update(arg1:mender.DeploymentStatus,arg2:number):Promise<void>;

export function WifiAPList():Promise<Array<string>>;

export function WifiConnect(arg1:string,arg2:string):Promise<void>;

export function WifiIsConnected():Promise<backend.WifiConnected>;
