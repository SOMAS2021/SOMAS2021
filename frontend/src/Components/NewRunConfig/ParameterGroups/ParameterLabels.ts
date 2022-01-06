import { SimConfig } from "../../../Helpers/SimConfig";

export interface Parameter {
  helperText: string;
  label: string;
  labelFor: string;
  labelInfo: string;
  key: keyof(SimConfig);
  min: number;
}