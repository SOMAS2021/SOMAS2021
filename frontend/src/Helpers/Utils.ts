import { DeathLog } from "./Logging/Death";

export function Average(arr: number[]): number {
  return arr.length > 0 ? arr.reduce((a, b) => a + b) / arr.length : 0;
}

// Deaths per agent
export function DeathsPerAgent(deathLog: DeathLog[]): { [agentType: string]: number } {
  var deaths: { [agentType: string]: number } = {};
  deathLog.forEach((death) => {
    deaths[death.agentType] = !deaths[death.agentType] ? 1 : deaths[death.agentType] + 1;
  });
  return deaths;
}