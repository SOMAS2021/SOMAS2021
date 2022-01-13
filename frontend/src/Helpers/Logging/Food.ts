import { GetFile } from "../API";
import { Log } from "./Log";

export interface FoodLog extends Log {
  food: number
}

export function GetFoodLogs(filename: string): Promise<FoodLog[]> {
  return new Promise<FoodLog[]>((resolve, reject) => {
    GetFile(filename, "food")
      .then((food) =>
        resolve(
          food.map(function (e: any) {
            const f: FoodLog = {
              food: e["food"],
              tick: e["tick"],
              day: e["day"],
            };
            return f;
          })
        )
      )
      .catch((err) => reject(err));
  });
}
