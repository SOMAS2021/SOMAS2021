import { Colors } from "@blueprintjs/core";
import { DeathLog } from "./Logging/Death";
import { UtilityLog } from "./Logging/Utility";
import { Result } from "./Result";

const arbitraryStackKey = "stack1";
const colorFamily = [
  "#1F4B99",
  "#2F649C",
  "#447C9F",
  "#5D94A1",
  "#7CAAA2",
  "#A0BFA2",
  "#CCD3A1",
  "#FFE39F",
  "#F6C880",
  "#EBAE65",
  "#DE944D",
  "#D07A38",
  "#C06126",
  "#B04718",
  "#9E2B0E",
];
const agentList = [
  "Team1Agents",
  "Team2Agents",
  "Team3Agents",
  "Team4Agents",
  "Team5Agents",
  "Team6Agents",
  "Team7Agents",
  "RandomAgents",
];

export function Average(arr: number[]): number {
  return arr.length > 0 ? arr.reduce((a, b) => a + b) / arr.length : 0;
}

export function Max(arr: number[]): number {
  return arr.length > 0 ? Math.max(...arr) : 0;
}

export function Min(arr: number[]): number {
  return arr.length > 0 ? Math.min(...arr) : 0;
}

// Deaths per agent
export function DeathsPerAgent(deathLog: DeathLog[]): { [agentType: string]: number } {
  var deaths: { [agentType: string]: number } = {};
  deathLog.forEach((death) => {
    deaths[death.agentType] = !deaths[death.agentType] ? 1 : deaths[death.agentType] + 1;
  });
  return deaths;
}

export function AverageUtilityPerAgent(utilityLogs: UtilityLog[]): { [agentType: string]: number } {
  // Return: the average utility per agent over its entire existence
  var utilities: { [agentType: string]: number } = {};
  var counts: { [agentType: string]: number } = {};
  // Sum utility and count num entries
  utilityLogs.forEach((log) => {
    let agentType = log.agentType;
    utilities[agentType] = !utilities[agentType] ? log.utility : utilities[agentType] + log.utility;
    counts[agentType] = !counts[agentType] ? 1 : counts[agentType] + 1;
  });

  // Get averages
  for (var agentType in utilities) {
    if (counts.hasOwnProperty(agentType)) {
      utilities[agentType] /= counts[agentType];
    }
  }
  return utilities;
}

export function UtilityOnDeath(utilities: UtilityLog[]): UtilityLog[] {
  return utilities.filter((o) => {
    return !o.isAlive;
  });
}

export function AverageAgeUponDeath(deathLogs: DeathLog[]): { [agentType: string]: number } {
  var ageMap: { [agentType: string]: number } = {};
  var counts: { [agentType: string]: number } = {};
  deathLogs.forEach((log) => {
    let agentType = log.agentType;
    ageMap[agentType] = !ageMap[agentType] ? log.ageUponDeath : ageMap[agentType] + log.ageUponDeath;
    counts[agentType] = !counts[agentType] ? 1 : counts[agentType] + 1;
  });
  // Get averages
  for (var agentType in ageMap) {
    if (counts.hasOwnProperty(agentType)) {
      ageMap[agentType] /= counts[agentType];
    }
  }
  return ageMap;
}

export function ParseMessageStats(result: Result): [string[], any[]] {
  let agentsPresent = GetPresentAgents(result);
  var newLabels: string[] = FilterArrayByOther(result.messages.atypes, agentsPresent);
  var newValues: any[] = [];
  for (let i = 0; i < result.messages.mtypes.length; i++) {
    newValues.push({
      label: result.messages.mtypes[i], // Graph title
      data: FilterArrayByOther(result.messages.msgcount[i], agentsPresent), // Data (y-axis)
      backgroundColor: colorFamily[i],
      stack: arbitraryStackKey,
    });
  }
  return [newLabels, newValues];
}

export function ParseTreatyAcceptanceStats(result: Result): [string[], any[]] {
  let agentsPresent = GetPresentAgents(result);
  let newLabels: string[] = FilterArrayByOther(result.messages.atypes, agentsPresent);
  let newValues: any[] = [
    {
      label: "Treaties Rejected",
      data: FilterArrayByOther(result.messages.treatyResponses[0], agentsPresent), // Data (y-axis)
      backgroundColor: Colors.RED1,
    },
    {
      label: "Treaties Accepted",
      data: FilterArrayByOther(result.messages.treatyResponses[1], agentsPresent), // Data (y-axis)
      backgroundColor: Colors.GREEN1,
    },
  ];
  return [newLabels, newValues];
}

export function ParseTreatyProposalStats(result: Result): [string[], any[]] {
  let agentsPresent = GetPresentAgents(result);
  var newLabels: string[] = FilterArrayByOther(result.messages.atypes, agentsPresent);
  var newValues: any[] = [];
  const i = 9;
  newValues.push({
    label: result.messages.mtypes[i], // Graph title
    data: FilterArrayByOther(result.messages.msgcount[i], agentsPresent), // Data (y-axis)
    backgroundColor: colorFamily[i],
    stack: arbitraryStackKey,
  });
  return [newLabels, newValues];
}

// Filter an array by an equal lengthed boolean array, no length check atm
export function FilterArrayByOther(data: any[], present: boolean[]): any[] {
  return data.filter((_, i) => {
    return present[i];
  });
}

export function GetPresentAgents(result: Result): boolean[] {
  var present: boolean[] = [];
  var config = result.config;
  agentList.forEach((agentType) => {
    const keyTyped = agentType as keyof typeof config;
    present.push(config[keyTyped] > 0);
  });
  return present;
}
