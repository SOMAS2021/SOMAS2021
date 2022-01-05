import { H3, H6, Divider, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetResult } from "../../Helpers/API";
import { Result } from "../../Helpers/Result";
import StoryViewer from "../Story/StoryViewer";
import { Average, Max, Min } from "../../Helpers/Utils";
import ConfigInfo from "./ConfigInfo";
import ReportCard from "./ReportCard";

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
    <div>
      {loading && <Spinner intent="primary" />}
      {!loading &&
        (logName !== "" && result ? (
          <ResultDisplay result={result} />
        ) : (
          <H6 style={{ paddingTop: 20 }}>
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
  return (
    <div
      style={{
        overflowY: "scroll",
        overflowX: "hidden",
        height: "95vh",
        textAlign: "left",
        padding: "20px 10px 30px 10px",
      }}
    >
      <H3>{result.title}</H3>
      <div>
        <Divider></Divider>
        <ConfigInfo config={result.config} />
        <Divider></Divider>
        <StoryViewer story={result.story} />
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
          <div className="col-lg-6">
            <ReportCard
              description="Average agent age upon death"
              title={Average(result.deaths.map((d) => d.ageUponDeath))
                .toFixed(3)
                .toString()}
            />
          </div>
          <div className="col-lg-6">
            <ReportCard
              description="Max agent age upon death"
              title={Max(result.deaths.map((d) => d.ageUponDeath))
                .toString()}
            />
          </div>
          <div className="col-lg-6">
            <ReportCard
              description="Min agent age upon death"
              title={Min(result.deaths.map((d) => d.ageUponDeath))
                .toString()}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
