import { H3, H5, H6, Spinner, Colors, Icon, Button, Divider } from "@blueprintjs/core";
import { useCallback, useEffect, useState } from "react";
import { GetResult, Result, SimStatusExec } from "../../Helpers/Result";
import CollapsingSection from "../CollapsingSection";
import StoryViewer from "../Story/StoryViewer";
import ConfigInfo from "./ConfigInfo";
import StatsViewer from "./Stats";

interface ResultsProps {
  logName: string;
}

export default function Results(props: ResultsProps) {
  const { logName } = props;
  const [result, setResult] = useState<Result>();
  const [loading, setLoading] = useState(false);

  const LoadResult = useCallback(() => {
    if (logName === "") return;
    setLoading(true);
    GetResult(logName)
      .then((result) => {
        setResult(result);
        setLoading(false);
      })
      .catch((_) => setLoading(false));
  }, [logName]);

  useEffect(() => {
    LoadResult();
  }, [logName, LoadResult]);

  return (
    <div>
      {loading && <Spinner intent="primary" />}
      {!loading &&
        (logName !== "" && result ? (
          <ResultDisplay result={result} reload={() => LoadResult()} />
        ) : (
          <H6 style={{ paddingTop: 20 }}>
            <i>Select an existing simulation from the sidebar to view results</i>
          </H6>
        ))}
    </div>
  );
}

interface ResultDisplayProps {
  result: Result;
  reload: () => void;
}
function ResultDisplay(props: ResultDisplayProps) {
  const { result, reload } = props;
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
      <ResultHeader result={result} reload={reload} />
      <Divider />
      <div>
        <CollapsingSection title="Config" defaultOpen={true}>
          <ConfigInfo config={result.config} initial={false} />
        </CollapsingSection>
        <StatsViewer result={result} />
        <CollapsingSection title="Story">
          <StoryViewer result={result} />
        </CollapsingSection>
      </div>
    </div>
  );
}

function ResultHeader(props: ResultDisplayProps) {
  const { result, reload } = props;
  return (
    <div className="row">
      <div className="col-lg-6">
        <H3>{result.title}</H3>
      </div>
      <div className="col-lg-6" style={{ textAlign: "right", paddingRight: "20px", height: 40 }}>
        {result.simStatus.status === SimStatusExec.timedout && (
          <H5 style={{ color: Colors.ORANGE2 }}>
            <Icon icon="warning-sign" size={20} style={{ paddingRight: 10 }} intent="warning" />
            Simulation timed out!
          </H5>
        )}
        {result.simStatus.status === SimStatusExec.finished && (
          <H5 style={{ color: Colors.GREEN2 }}>
            <Icon icon="tick" size={20} style={{ paddingRight: 10 }} intent="success" />
            Simulation completed!
          </H5>
        )}
        {result.simStatus.status === SimStatusExec.running && (
          <>
            <H5 style={{ color: Colors.BLUE2 }}>
              <Icon className="rotate-image" icon="refresh" size={20} style={{ padding: "0 10px" }} intent="primary" />
              Simulation still running!
              <Button text="Refresh" intent="primary" style={{ marginLeft: 10 }} onClick={reload} />
            </H5>
          </>
        )}
      </div>
    </div>
  );
}
