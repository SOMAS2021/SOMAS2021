import { GetFile } from "../API";
import { Log } from "./Log";

export interface UtilityLog extends Log {
  agentType: string;
  utility: number;
  isAlive: boolean;
}

export function GetUtilityLogs(filename: string): Promise<UtilityLog[]> {
  return new Promise<UtilityLog[]>((resolve, reject) => {
    GetFile(filename, "utility")
      .then((u) =>
        resolve(
          u.map(function (e: any) {
            const d: UtilityLog = {
              agentType: e["agent_type"],
              tick: e["tick"],
              day: e["day"],
              utility: e["utility"],
              isAlive: e["isAlive"],
            };
            return d;
          })
        )
      )
      .catch((err) => reject(err));
  });
}
