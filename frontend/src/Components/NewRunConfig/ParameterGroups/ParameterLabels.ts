import { SimConfig } from "../../../Helpers/SimConfig";

export interface Parameter {
  helperText: string;
  label: string;
  labelFor: string;
  labelInfo: string;
  key: keyof(SimConfig);
  min: number;
}

// might not need these
// export const healthParams: Parameter[] = [
//   {
//     helperText: "",
//     label: "Weak Level",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "weakLevel",
//     min:1
//   },
//   {
//     helperText: "",
//     label: "Width",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "width",
//     min:1
//   },
//   {
//     helperText: "",
//     label: "Tau",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "tau",
//     min:1
//   },
//   {
//     helperText: "Upper bound",
//     label: "Critical HP",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "hpCritical",
//     min:1
//   },
//   {
//     helperText: "",
//     label: "Days at Critical",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "maxDayCritical",
//     min:1
//   },
//   {
//     helperText: "",
//     label: "HP Loss Base",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "HPLossBase",
//     min:1
//   },
//   {
//     helperText: "",
//     label: "HP Loss Slope",
//     labelFor: "text-input",
//     labelInfo: "",
//     key: "HPLossSlope",
//     min:0
//   },
// ];