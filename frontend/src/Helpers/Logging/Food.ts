import { GetFile } from "../API";
import { Log } from "./Log";

export interface FoodLog extends Log {
  food: number;
}
export interface FoodFloorLog extends Log {
  food: number;
  floor: number;
}

export function GetFoodLogs(filename: string): Promise<FoodLog[]> {
  return new Promise<FoodLog[]>((resolve, reject) => {
    GetFile(filename, "foodDay")
      .then((food) => {
        const foods: FoodLog[] = food.map(function (e: any) {
          const f: FoodLog = {
            food: e["food"],
            tick: e["tick"],
            day: e["day"],
          };
          return f;
        });
        resolve(foods);
      })
      .catch((err) => reject(err));
  });
}

export function GetFoodFloorLogs(filename: string): Promise<FoodFloorLog[]> {
  return new Promise<FoodFloorLog[]>((resolve, reject) => {
    GetFile(filename, "foodFloor")
      .then((food) => {
        const foods: FoodFloorLog[] = food.map(function (e: any) {
          const f: FoodFloorLog = {
            food: e["food"],
            tick: e["tick"],
            day: e["day"],
            floor: e["floor"],
          };
          return f;
        });
        resolve(foods);
      })
      .catch((err) => reject(err));
  });
}
