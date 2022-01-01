export const advancedParams: Parameters[] = [
  {
    helperText: "",
    label: "Agents Max HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentHP",
  },
  {
    helperText: "",
    label: "Agents Per Floor",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentsPerFloor",
  },
  {
    helperText: "In Days...",
    label: "Reshuffle Period",
    labelFor: "text-input",
    labelInfo: "",
    key: "ReshuffleDays",
  },
  {
    helperText: "",
    label: "Weak Level",
    labelFor: "text-input",
    labelInfo: "",
    key: "weakLevel",
  },
  {
    helperText: "",
    label: "Width",
    labelFor: "text-input",
    labelInfo: "",
    key: "width",
  },
  {
    helperText: "",
    label: "Tau",
    labelFor: "text-input",
    labelInfo: "",
    key: "tau",
  },
  {
    helperText: "Upper bound",
    label: "Critical HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "hpCritical",
  },
  {
    helperText: "",
    label: "Days at Critical",
    labelFor: "text-input",
    labelInfo: "",
    key: "maxDayCritical",
  },
  {
    helperText: "",
    label: "HP Loss Base",
    labelFor: "text-input",
    labelInfo: "",
    key: "HPLossBase",
  },
  {
    helperText: "",
    label: "HP Loss Slope",
    labelFor: "text-input",
    labelInfo: "",
    key: "HPLossSlope",
  },
];

export const params: Parameters[] = [
    {
      helperText: "In Days...",
      label: "Simulation Length",
      labelFor: "text-input",
      labelInfo: "",
      key: "SimDays",
    },
    {
      helperText: "Agents can do one/two actions per 'Tick'",
      label: "'Ticks' Per Floor",
      labelFor: "text-input",
      labelInfo: "",
      key: "TicksPerFloor",
    },
    {
      helperText: "Food on the platfrom at the beginnning of each day",
      label: "Initial Food On The Platform",
      labelFor: "text-input",
      labelInfo: "",
      key: "FoodOnPlatform",
    },
    {
      helperText: "Number of Agents from Team 1 in the Tower",
      label: "Agents of Team 1",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team1Agents",
    },
    {
      helperText: "Number of Agents from Team 2 in the Tower",
      label: "Agents of Team 2",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team2Agents",
    },
    {
      helperText: "Number of Agents from Team 3 in the Tower",
      label: "Agents of Team 3",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team3Agents",
    },
    {
      helperText: "Number of Agents from Team 4 in the Tower",
      label: "Agents of Team 4",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team4Agents",
    },
    {
      helperText: "Number of Agents from Team 5 in the Tower",
      label: "Agents of Team 5",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team5Agents",
    },
    {
      helperText: "Number of Agents from Team 6 in the Tower",
      label: "Agents of Team 6",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team6Agents",
    },
    {
      helperText: "Number of Agents from Team 7 in the Tower",
      label: "Agents of Team 7",
      labelFor: "text-input",
      labelInfo: "",
      key: "Team7Agents",
    },
    {
      helperText: "Number of Random Agents in the Tower",
      label: "Random Agents",
      labelFor: "text-input",
      labelInfo: "",
      key: "RandomAgents",
    },
];

export interface Parameters {
  helperText: string;
  label: string;
  labelFor: string;
  labelInfo: string;
  key: string;
}
