import { GetFile } from "./API";
import { UtilityLog } from "./Logging/Utility";
import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { GetStoryLogs, StoryLog } from "./Logging/StoryLog";
import { GetSimConfig, SimConfig } from "./SimConfig";
import { GetUtilityLogs } from "./Logging/Utility";

export enum SimStatus {
  "finished",
  "running",
  "timedout",
}

export interface Result {
  title: string;
  deaths: DeathLog[];
  food: FoodLog[];
  config: SimConfig;
  story: StoryLog[];
  status: SimStatus;
  utility: UtilityLog[];
}

function GetSimStatus(filename: string): Promise<SimStatus> {
  return new Promise<SimStatus>((resolve, reject) => {
    GetFile(filename, "status")
      .then((status) => {
        switch (status[0]["status"]) {
          case SimStatus[SimStatus.finished]:
            resolve(SimStatus.finished);
            break;
          case SimStatus[SimStatus.timedout]:
            resolve(SimStatus.timedout);
            break;
          case SimStatus[SimStatus.running]:
            resolve(SimStatus.running);
            break;
          default:
            reject("unkown status");
        }
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
    var status: SimStatus = SimStatus.running;
    promises.push(GetSimStatus(filename).then((s) => (status = s)));

    // story
    var story: StoryLog[] = [];
    promises.push(GetStoryLogs(filename).then((s) => (story = s)));

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
        story: story,
        status: status,
        utility: utility,
      })
    );
  });
}
