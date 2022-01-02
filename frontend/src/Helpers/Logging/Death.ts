import { GetFile } from "../API";
import { Log } from "./Log";

export interface DeathLog extends Log{
  cumulativeDeaths: number
  agentType: string
}

export function GetDeathLogs(filename: string): Promise<DeathLog[]> {
  return new Promise<DeathLog[]>((resolve, reject) => {
    GetFile(filename, "death")
      .then((deaths) =>
        resolve(
          deaths.map(function (e: any) {
            const d: DeathLog = {
              cumulativeDeaths: e["cumulativeDeaths"],
              agentType: e["agent_type"],
              tick: e["tick"],
              day: e["day"],
            };
            return d;
          })
        )
      )
      .catch((err) => reject(err));
  });
}
