import { DeathLog } from "./Logging/Death";
import { FoodLog } from "./Logging/Food";

export interface Result {
  title: string
  deaths: DeathLog[]
  food: FoodLog[]
}