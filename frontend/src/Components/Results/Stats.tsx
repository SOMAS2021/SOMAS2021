import { Result, SimStatusExec } from "../../Helpers/Result";
import CollapsingSection from "../CollapsingSection";
import DeathStats from "./Stats/DeathStats";
import OtherStats from "./Stats/OtherStats";
export interface StatsViewerProps {
  result: Result;
}

export default function StatsViewer(props: StatsViewerProps) {
  const { result } = props;
  const disabled = result.simStatus.status === SimStatusExec.running;
  return (
    <div>
      <CollapsingSection title="Death Stats" defaultOpen={true}>
        <DeathStats result={result} />
      </CollapsingSection>
      <CollapsingSection title="Other Stats" disabled={disabled}>
        <OtherStats result={result} />
      </CollapsingSection>
    </div>
  );
}
