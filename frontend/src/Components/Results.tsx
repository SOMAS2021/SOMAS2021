import { H3, H6, Divider, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetResult } from "../Helpers/API";
import { Result } from "../Helpers/Result";
import { Average } from "../Helpers/Utils";
import ReportCard from "./Results/ReportCard";

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
  return (
    <>
      <H3>{result.title}</H3>
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
    </>
  );
}
