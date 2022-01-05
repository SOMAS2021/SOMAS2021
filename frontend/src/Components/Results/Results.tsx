import { H3, H6, Divider, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetResult } from "../../Helpers/API";
import { Result } from "../../Helpers/Result";
import { Average } from "../../Helpers/Utils";
import ConfigInfo from "./ConfigInfo";
import ReportCard from "./ReportCard";
import LineChart from "./Graphs/LineChart";
import BarChart from "./Graphs/BarChart";
interface ResultsProps {
  logName: string;
}

export default function Results(props: ResultsProps) {
  const { logName } = props;
  const [result, setResult] = useState<Result>();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    if (logName === "") return;
    setLoading(true);
    GetResult(logName)
      .then((result) => {
        setResult(result);
        setLoading(false);
      })
      .catch((_) => setLoading(false));
  }, [logName]);

  return (
    <div style={{ padding: 20 }}>
      {loading && <Spinner intent="primary" />}
      {!loading &&
        (logName !== "" && result ? (
          <ResultDisplay result={result} />
        ) : (
          <H6>
            <i>Select an existing simulation result to view results</i>
          </H6>
        ))}
    </div>
  );
}

interface ResultDisplayProps {
  result: Result;
}
function ResultDisplay(props: ResultDisplayProps) {
  const { result } = props;

// CUMULATIVE DEATHS PER AGENT CODE
// Finding which agent types are duplicated and their respective indices
// Adding the cumulative deaths, and removing duplicated agentTypes
  var agentTypes = result.deaths.map((d) => d.agentType) 
  var cumulativeDeaths = result.deaths.map((d) => d.cumulativeDeaths) 
  var myMap = new Map<string, number[]>()


  for (var i = 0; i < agentTypes.length; i++) {
      var element = agentTypes[i].toString();  
      var test = myMap.get(element)
      if (typeof test != "undefined"){
        test.push(i)
        myMap.set(element,test)
      }else {
        myMap.set(element,[i])
      }
  }

  var findDuplicates = Array.from(myMap.entries()).filter(
    ([key, val]) =>{
      if (val.length > 1){
        return [key, val]
      }
    });
    
  if (findDuplicates){
    for (var j=0; j<findDuplicates.length; j++){
      let [n, idx] = findDuplicates[j]
      var tempValue = 0
      for(var i = 0; i<idx.length; i++){
        tempValue += cumulativeDeaths[idx[i]]
      }
      cumulativeDeaths[idx[0]] = tempValue

      for (var z = 1; z<idx.length; z++){
        delete cumulativeDeaths[idx[z]]
        delete agentTypes[idx[z]]
      }
    }
    cumulativeDeaths = cumulativeDeaths.filter(function(e){return e});
    agentTypes = agentTypes.filter(function(e){return e});

  }

  return (
    <>
      <H3>{result.title}</H3>
      <Divider></Divider>
      <ConfigInfo config={result.config} />
      <Divider></Divider>
      <div className="row">
        <div className="col-lg-6">
          <ReportCard description="Total deaths" title={result.deaths.length.toString()} />
        </div>
        <div className="col-lg-6">
          <ReportCard
            description="Average food on platform per tick"
            title={Average(result.food.map((f) => f.food))
              .toFixed(3)
              .toString()}
          />
        </div>
      </div>
      <div className="row">
        <div className="col-lg-6">
          <BarChart yAxis={cumulativeDeaths} xAxis={agentTypes} graphTitle="Cumulative Deaths per Agent type" />
        </div>
        <div className="col-lg-6">
          <LineChart yAxis={result.food.map((f)=> f.food)} xAxis={result.food.map((f) => f.tick)} graphTitle="Total Food on Platform per Tick"/>
        </div>
      </div>
    </>
  );
}
