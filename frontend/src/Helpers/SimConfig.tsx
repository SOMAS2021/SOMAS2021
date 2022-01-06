import { useState } from "react";
import { GetFile } from "./API";

export default function InitConfigState() {
  return useState<SimConfig>({
    FoodOnPlatform: 100,
    FoodPerAgentRatio: 10,
    UseFoodPerAgentRatio: false,
    Team1Agents: 2,
    Team2Agents: 2,
    Team3Agents: 2,
    Team4Agents: 2,
    Team5Agents: 2,
    Team6Agents: 2,
    Team7Agents: 2,
    RandomAgents: 2,
    AgentHP: 100,
    AgentsPerFloor: 1,
    TicksPerFloor: 10,
    SimDays: 8,
    ReshuffleDays: 1,
    maxHP: 100,
    weakLevel: 10,
    width: 48,
    tau: 15,
    hpReqCToW: 5,
    hpCritical: 3,
    maxDayCritical: 3,
    HPLossBase: 5,
    HPLossSlope: 0.2,
    LogFileName: "",
    LogMain: false,
    SimTimeoutSeconds: 10,
  });
}

export interface SimConfig {
  FoodOnPlatform: number;
  FoodPerAgentRatio: number;
  UseFoodPerAgentRatio: boolean;
  Team1Agents: number;
  Team2Agents: number;
  Team3Agents: number;
  Team4Agents: number;
  Team5Agents: number;
  Team6Agents: number;
  Team7Agents: number;
  RandomAgents: number;
  AgentHP: number;
  AgentsPerFloor: number;
  TicksPerFloor: number;
  SimDays: number;
  ReshuffleDays: number;
  maxHP: number;
  weakLevel: number;
  width: number;
  tau: number;
  hpReqCToW: number;
  hpCritical: number;
  maxDayCritical: number;
  HPLossBase: number;
  HPLossSlope: number;
  LogFileName: string;
  LogMain: boolean;
  SimTimeoutSeconds: number;
}

export function GetSimConfig(filename: string): Promise<SimConfig> {
  return new Promise<SimConfig>((resolve, reject) => {
    GetFile(filename, "config")
      .then((config) => {
        resolve(config[0] as SimConfig);
      })
      .catch((err) => reject(err));
  });
}
