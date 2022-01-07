import { GetFile } from "./API";
import { AgentStateLog } from "./Logging/AgentState";
import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { GetStoryLogs, StoryLog } from "./Logging/StoryLog";
import { GetSimConfig, SimConfig } from "./SimConfig";
import { GetAgentStateLogs } from "./Logging/AgentState";

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
  agents: AgentStateLog[];
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
    var agents: AgentStateLog[] = [];
    promises.push(GetAgentStateLogs(filename).then((a) => (agents = a)));

    // all
    Promise.all(promises).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
        food: foods,
        config: config,
        story: story,
        status: status,
        agents: agents,
      })
    );
  });
}
