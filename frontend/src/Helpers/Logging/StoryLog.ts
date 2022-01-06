import { GetFile } from "../API";
import { Log } from "./Log";

export interface StoryLog extends Log {
  msg: string;
}

export interface AgentState {
  hp: number;
  atype: string;
  age: number;
  floor: number;
  state: string;
}

export interface StoryFoodLog extends StoryLog, AgentState {
  foodTaken: number;
  foodLeft: number;
}

export interface StoryMessageLog extends StoryLog, AgentState {
  target: number;
  mtype: string;
  mcontent: string;
}

export interface StoryDeathLog extends StoryLog, AgentState {}

export interface StoryPlatformLog extends StoryLog {
  floor: number;
}

export function GetStoryLogs(filename: string): Promise<StoryLog[]> {
  return new Promise<StoryLog[]>((resolve, reject) => {
    GetFile(filename, "story")
      .then((storylogs) =>
        resolve(
          storylogs.map(function (e: any) {
            return e as StoryLog;
          })
        )
      )
      .catch((err) => reject(err));
  });
}




export function GetStoryMessageLogs(filename: string): Promise<StoryMessageLog[]> {
  return new Promise<StoryMessageLog[]>((resolve, reject) => {
    GetFile(filename, "story")
      .then((storymessagelogs) =>
        resolve(
          storymessagelogs.map(function (e: any) {
            const m: StoryMessageLog = {
              target: e["target"],
              mtype: e["mtype"],
              mcontent: e["mcontent"],
              msg: e["msg"],
              hp: e["hp"],
              atype: e["atype"],
              age: e["age"],
              floor: e["floor"], 
              state: e["state"],
              tick: e["tick"],
              day: e["day"],
            };
            return m;
          })
        )
      )
      .catch((err) => reject(err));
  });
}

