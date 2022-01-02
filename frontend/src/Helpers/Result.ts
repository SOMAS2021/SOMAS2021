import { DeathLog } from "./Logs/Death";
import { FoodLog } from "./Logs/Food";

export interface Result {
  title: string
  deaths: DeathLog[]
  food: FoodLog[]
}