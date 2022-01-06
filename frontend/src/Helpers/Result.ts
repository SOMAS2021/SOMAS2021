import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { GetStoryLogs, StoryLog } from "./Logging/StoryLog";
import { GetSimConfig, SimConfig } from "./SimConfig";

export interface Result {
  title: string
  deaths: DeathLog[]
  food: FoodLog[]
  config: SimConfig
  story: StoryLog[]
}

export function GetResult(filename: string): Promise<Result> {
  return new Promise<Result>(async (resolve, reject) => {
    // promises
    var promises: Promise<any>[] = [];
    // deaths
    var deaths: DeathLog[] = [];
    promises.push(GetDeathLogs(filename).then((d) => (deaths = d)));

    // foods
    var foods: FoodLog[] = [];
    promises.push(GetFoodLogs(filename).then((f) => (foods = f)));

    // config
    var config: SimConfig = undefined!;
    promises.push(GetSimConfig(filename).then((c) => (config = c)));

    // story
    var story: StoryLog[] = [];
    promises.push(GetStoryLogs(filename).then((s) => (story = s)));

    // all
    Promise.all(promises).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
        food: foods,
        config: config,
        story: story,
      })
    );
  });
}