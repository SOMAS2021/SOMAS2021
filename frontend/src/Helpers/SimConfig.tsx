import { useState } from "react";

export default function InitConfigState() {
  return useState<SimConfig>({
    FoodOnPlatform: 100,
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
    width: 45,
    tau: 10,
    hpReqCToW: 2,
    hpCritical: 5,
    maxDayCritical: 3,
    HPLossBase: 10,
    HPLossSlope: 0.25,
    LogFileName:"", 
    LogMain: false,
    SimTimeoutSeconds: 10
  });
}

export interface SimConfig {
  FoodOnPlatform: number;
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
  LogFileName:string, 
  LogMain: boolean,
  SimTimeoutSeconds: number
}
