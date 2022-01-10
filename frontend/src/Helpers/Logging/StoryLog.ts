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

export function GetStoryLogs(filename: string, tick: number): Promise<StoryLog[]> {
  return new Promise<StoryLog[]>((resolve, reject) => {
    GetFile(filename, "story", tick)
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
