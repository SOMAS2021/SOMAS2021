import { DeathLog } from "./Logging/Death";
import { UtilityLog } from "./Logging/Utility";

export function Average(arr: number[]): number {
  return arr.length > 0 ? arr.reduce((a, b) => a + b) / arr.length : 0;
}

export function Max(arr: number[]): number {
  return arr.length > 0 ? Math.max(...arr) : 0;
}

export function Min(arr: number[]): number {
  return arr.length > 0 ? Math.min(...arr) : 0;
}

// Deaths per agent
export function DeathsPerAgent(deathLog: DeathLog[]): { [agentType: string]: number } {
  var deaths: { [agentType: string]: number } = {};
  deathLog.forEach((death) => {
    deaths[death.agentType] = !deaths[death.agentType] ? 1 : deaths[death.agentType] + 1;
  });
  return deaths;
}

export function AverageUtilityPerAgent(utilityLogs: UtilityLog[]): { [agentType: string]: number } {
  // Return: the average utility per agent over its entire existence
  var utilities: { [agentType: string]: number } = {};
  var counts: { [agentType: string]: number } = {};
  // Sum utility and count num entries
  utilityLogs.forEach((log) => {
    let agentType = log.agentType;
    utilities[agentType] = !utilities[agentType] ? log.utility : utilities[agentType] + log.utility;
    counts[agentType] = !counts[agentType] ? 1 : counts[agentType] + 1;
  });

  // Get averages
  for (var agentType in utilities) {
    if (counts.hasOwnProperty(agentType)) {
      utilities[agentType] /= counts[agentType];
    }
  }
  return utilities;
}

export function UtilityOnDeath(utilities: UtilityLog[]): UtilityLog[] {
  return utilities.filter((o) => {
    return !o.isAlive;
  });
}

export function AverageAgeUponDeath(deathLogs: DeathLog[]): { [agentType: string]: number } {
  var ageMap: { [agentType: string]: number } = {};
  var counts: { [agentType: string]: number } = {};
  deathLogs.forEach((log) => {
    let agentType = log.agentType;
    ageMap[agentType] = !ageMap[agentType] ? log.ageUponDeath : ageMap[agentType] + log.ageUponDeath;
    counts[agentType] = !counts[agentType] ? 1 : counts[agentType] + 1;
  });
  // Get averages
  for (var agentType in ageMap) {
    if (counts.hasOwnProperty(agentType)) {
      ageMap[agentType] /= counts[agentType];
    }
  }
  return ageMap;
}
