import { SimConfig } from "../../Helpers/SimConfig";

export const advancedParams: Parameter[] = [
  {
    helperText: "",
    label: "Agents Max HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentHP",
    min:1
  },
  {
    helperText: "",
    label: "Agents Per Floor",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentsPerFloor",
    min:1
  },
  {
    helperText: "In Days...",
    label: "Reshuffle Period",
    labelFor: "text-input",
    labelInfo: "",
    key: "ReshuffleDays",
    min:1
  },
  {
    helperText: "",
    label: "Weak Level",
    labelFor: "text-input",
    labelInfo: "",
    key: "weakLevel",
    min:1
  },
  {
    helperText: "",
    label: "Width",
    labelFor: "text-input",
    labelInfo: "",
    key: "width",
    min:1
  },
  {
    helperText: "",
    label: "Tau",
    labelFor: "text-input",
    labelInfo: "",
    key: "tau",
    min:1
  },
  {
    helperText: "Upper bound",
    label: "Critical HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "hpCritical",
    min:1
  },
  {
    helperText: "",
    label: "Days at Critical",
    labelFor: "text-input",
    labelInfo: "",
    key: "maxDayCritical",
    min:1
  },
  {
    helperText: "",
    label: "HP Loss Base",
    labelFor: "text-input",
    labelInfo: "",
    key: "HPLossBase",
    min:1
  },
  {
    helperText: "",
    label: "HP Loss Slope",
    labelFor: "text-input",
    labelInfo: "",
    key: "HPLossSlope",
    min:0
  },
];

export const params: Parameter[] = [
    {
      helperText: "In Days...",
      label: "Simulation Length",
      labelFor: "text-input",
      labelInfo: "",
      key: "SimDays",
      min:1
    },
    {
      helperText: "Agents can do one/two actions per 'Tick'",
      label: "'Ticks' Per Floor",
      labelFor: "text-input",
      labelInfo: "",
      key: "TicksPerFloor",
      min:1
    },
    {
      helperText: "Food on the platform at the beginning of each day",
      label: "Initial Food On The Platform",
      labelFor: "text-input",
      labelInfo: "",
      key: "FoodOnPlatform",
      min:1,
    },
    {
      helperText: "Food on the platform at the beginning of each day",
      label: "Food Per Agent",
      labelFor: "text-input",
      labelInfo: "",
      key: "FoodPerAgentRatio",
      min:1,
    },
    {
      helperText: "Number of Agents from Team 1 in the Tower",
      label: "Agents of Team 1",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team1Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 2 in the Tower",
      label: "Agents of Team 2",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team2Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 3 in the Tower",
      label: "Agents of Team 3",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team3Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 4 in the Tower",
      label: "Agents of Team 4",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team4Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 5 in the Tower",
      label: "Agents of Team 5",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team5Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 6 in the Tower",
      label: "Agents of Team 6",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team6Agents",
      min:0,
    },
    {
      helperText: "Number of Agents from Team 7 in the Tower",
      label: "Agents of Team 7",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team7Agents",
      min:0,
    },
    {
      helperText: "Number of Random Agents in the Tower",
      label: "Random Agents",
      labelFor: "text-input",
      labelInfo: "",
      key: "RandomAgents",
      min:0,
    },
    {
      helperText: "Base this on the size of your simulation.",
      label: "TimeOut",
      labelFor: "text-input",
      labelInfo: "",
      key: "SimTimeoutSeconds",
      min:1,
    }
];

export interface Parameter {
  helperText: string;
  label: string;
  labelFor: string;
  labelInfo: string;
  key: keyof(SimConfig);
  min: number;
}
