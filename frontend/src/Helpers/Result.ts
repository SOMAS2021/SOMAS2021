import { DeathLog } from "./Logging/Death";
import { FoodLog } from "./Logging/Food";
import { SimConfig } from "./SimConfig";

export interface Result {
  title: string
  deaths: DeathLog[]
  food: FoodLog[]
  config: SimConfig
}