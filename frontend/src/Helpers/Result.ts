import { DeathLog } from "./Logging/Death";
import { FoodLog } from "./Logging/Food";
import { StoryLog } from "./Logging/StoryLog";
import { SimConfig } from "./SimConfig";

export interface Result {
  title: string
  deaths: DeathLog[]
  food: FoodLog[]
  config: SimConfig
  story: StoryLog[]
}