import { GetFile } from "../API";
import { Log } from "./Log";

export interface FoodLog extends Log {
  food: number;
}

export function GetFoodLogs(filename: string): Promise<FoodLog[]> {
  return new Promise<FoodLog[]>((resolve, reject) => {
    GetFile(filename, "food")
      .then((food) => {
        const foods: FoodLog[] = food.map(function (e: any) {
          const f: FoodLog = {
            food: e["food"],
            tick: e["tick"],
            day: e["day"],
          };
          return f;
        });
        // var maxTick: number = foods[foods.length - 1].tick;
        // for (var i = 0; i < maxTick; i++) {
        //   if (!foods[i] || foods[i].tick !== i + 1) {
        //     foods.splice(i, 0, { tick: i + 1, day: foods[i - 1].day, food: foods[i - 1].food });
        //   }
        // }
        resolve(foods);
      })
      .catch((err) => reject(err));
  });
}
