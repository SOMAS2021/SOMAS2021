import { GetFile } from "./API";
import { UtilityLog } from "./Logging/Utility";
import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { GetSimConfig, SimConfig } from "./SimConfig";
import { GetUtilityLogs } from "./Logging/Utility";

export enum SimStatusExec {
  "finished",
  "running",
  "timedout",
}

export interface SimStatus {
  status: SimStatusExec;
  maxTick: number;
}

export interface Result {
  title: string;
  deaths: DeathLog[];
  food: FoodLog[];
  config: SimConfig;
  simStatus: SimStatus;
  utility: UtilityLog[];
}

function GetSimStatus(filename: string): Promise<SimStatus> {
  return new Promise<SimStatus>((resolve, reject) => {
    GetFile(filename, "status")
      .then((status) => {
        var s: SimStatus = {
          status: SimStatusExec.finished,
          maxTick: status[0]["maxTick"],
        };
        switch (status[0]["status"]) {
          case SimStatusExec[SimStatusExec.finished]:
            s.status = SimStatusExec.finished;
            break;
          case SimStatusExec[SimStatusExec.timedout]:
            s.status = SimStatusExec.timedout;
            break;
          case SimStatusExec[SimStatusExec.running]:
            s.status = SimStatusExec.running;
            break;
          default:
            reject("unkown status");
        }
        resolve(s);
      })
      .catch((err) => reject(err));
  });
}

export function GetResult(filename: string): Promise<Result> {
  return new Promise<Result>((resolve, reject) => {
    // promises
    var promises: Promise<any>[] = [];
    // deaths
    var deaths: DeathLog[] = [];
    promises.push(GetDeathLogs(filename).then((d) => (deaths = d)));

    // foods
    var foods: FoodLog[] = [];
    promises.push(GetFoodLogs(filename).then((f) => (foods = f)));

    // config
    var config: SimConfig = undefined!;
    promises.push(GetSimConfig(filename).then((c) => (config = c)));

    // status
    var status: SimStatus = {
      status: SimStatusExec.finished,
      maxTick: 3000
    };
    promises.push(GetSimStatus(filename).then((s) => (status = s)));

    // Agent state
    var utility: UtilityLog[] = [];
    promises.push(GetUtilityLogs(filename).then((a) => (utility = a)));

    // all
    Promise.all(promises).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
        food: foods,
        config: config,
        simStatus: status,
        utility: utility,
      })
    );
  });
}
