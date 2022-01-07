import { GetFile } from "../API";
import { Log } from "./Log";

export interface AgentStateLog extends Log {
  agentType: string;
  utility: number;
  hp: number;
  floor: number;
}

export function GetAgentStateLogs(filename: string): Promise<AgentStateLog[]> {
  return new Promise<AgentStateLog[]>((resolve, reject) => {
    GetFile(filename, "AgentState")
      .then((AgentStates) =>
        resolve(
          AgentStates.map(function (e: any) {
            const d: AgentStateLog = {
              agentType: e["agent_type"],
              tick: e["tick"],
              day: e["day"],
              utility: e["utility"],
              hp: e["hp"],
              floor: e["floor"],
            };
            return d;
          })
        )
      )
      .catch((err) => reject(err));
  });
}
