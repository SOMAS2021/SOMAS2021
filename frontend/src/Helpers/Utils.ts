import { DeathLog } from "./Logging/Death";
import { StoryLog, StoryMessageLog } from "./Logging/StoryLog";

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

export function MessagesPerAgent(messageLog: StoryMessageLog[]) {
  var message: { [agentType: string]: number } = {};
  messageLog.forEach((m) => {
    message[m.atype] = !message[m.atype] ? 1 : message[m.atype] + 1;
  });
  return message;
}
