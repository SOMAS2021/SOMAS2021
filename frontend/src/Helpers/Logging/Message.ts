import { GetFile } from "../API";
import { Log } from "./Log";

export interface MessagesLog extends Log {
  atypes: string[];
  msgcount: number[][];
  mtypes: string[];
  treatyResponses: number[][];
}

export function GetMessagesLog(filename: string): Promise<MessagesLog> {
  return new Promise<MessagesLog>((resolve, reject) => {
    GetFile(filename, "messages")
      .then((u) => {
        resolve(u[0] as MessagesLog);
      })
      .catch((err) => reject(err));
  });
}
