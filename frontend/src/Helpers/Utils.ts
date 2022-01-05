import { DeathLog } from "./Logging/Death";

export function Average(arr: number[]): number {
  return arr.reduce((a, b) => a + b) / arr.length;
}

// CUMULATIVE DEATHS PER AGENT CODE
// Finding which agent types are duplicated and their respective indices
// Adding the cumulative deaths, and removing duplicated agentTypes
export function addCumulativeDeaths(deathLog: DeathLog[]) {
  var agentTypes = deathLog.map((d) => d.agentType);
  var cumulativeDeaths = deathLog.map((d) => d.cumulativeDeaths);
  var myMap = new Map<string, number[]>();

  for (var i = 0; i < agentTypes.length; i++) {
    var element = agentTypes[i].toString();
    var test = myMap.get(element);
    if (typeof test != "undefined") {
      test.push(i);
      myMap.set(element, test);
    } else {
      myMap.set(element, [i]);
    }
  }

  var findDuplicates = Array.from(myMap.entries()).filter(([key, val]) => {
    if (val.length > 1) {
      return [key, val];
    }
    return [];
  });

  if (findDuplicates) {
    for (var j = 0; j < findDuplicates.length; j++) {
      let idx = findDuplicates[j][1];
      var tempValue = 0;
      for (i = 0; i < idx.length; i++) {
        tempValue += cumulativeDeaths[idx[i]];
      }
      cumulativeDeaths[idx[0]] = tempValue;

      for (var z = 1; z < idx.length; z++) {
        delete cumulativeDeaths[idx[z]];
        delete agentTypes[idx[z]];
      }
    }
    cumulativeDeaths = cumulativeDeaths.filter(function (e) {
      return e;
    });
    agentTypes = agentTypes.filter(function (e) {
      return e;
    });
  }

  return { agentTypes, cumulativeDeaths };
}
