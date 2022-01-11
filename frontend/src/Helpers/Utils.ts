import { Colors } from "@blueprintjs/core";
import { DeathLog } from "./Logging/Death";
import { MessagesLog } from "./Logging/Message";
import { UtilityLog } from "./Logging/Utility";

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

// This needs to return an array of chartable objects
export function ParseMessageStats(stats: MessagesLog): any[] {
  var processedStats: any[] = [];
  var mtypes = stats.mtypes;
  var msgcount = stats.msgcount;

  // for each agent type
  for (let i = 0; i < mtypes.length; i++) {
    processedStats.push({
      label: mtypes[i], // Graph title
      data: msgcount[i], // Data (y-axis)
      backgroundColor: colorFamily[i],
      stack: arbitraryStackKey,
    });
  }

  return processedStats;
}
export function ParseTreatyAcceptanceStats(stats: MessagesLog): any[] {
  return [
    {
      label: "Treaties Rejected", // Graph title
      data: stats.treatyResponses[0], // Data (y-axis)
      backgroundColor: Colors.RED1,
    },
    {
      label: "Treaties Accepted", // Graph title
      data: stats.treatyResponses[1], // Data (y-axis)
      backgroundColor: Colors.GREEN3,
    },
  ];
}
