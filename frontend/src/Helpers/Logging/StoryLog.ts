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
