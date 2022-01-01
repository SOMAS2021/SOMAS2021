import { Card, Elevation, H3, H1, H6, Divider, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetResult } from "../Helpers/API";
import { Result } from "../Helpers/Result";

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
            <i>Select an exsiting simulation result to view results</i>
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
      <H3>{result!.title}</H3>
      <Divider></Divider>
      <div className="row">
        <div className="col-lg-2">
          <ReportCard description="Total deaths" title={result!.deaths.slice(-1)[0].cumulativeDeaths.toString()} />
        </div>
      </div>
    </>
  );
}

interface ReportCardProps {
  title: string;
  description: string;
}

function ReportCard(props: ReportCardProps) {
  const { title, description } = props;
  return (
    <Card interactive={true} elevation={Elevation.TWO} style={{ marginTop: 20 }}>
      <H1 style={{ color: "#1F4B99" }}>{title}</H1>
      <p>{description}</p>
    </Card>
  );
}
