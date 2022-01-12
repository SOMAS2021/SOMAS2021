import { Result } from "../../Helpers/Result";
import CollapsingSection from "../CollapsingSection";
import DeathStats from "./Stats/DeathStats";
import OtherStats from "./Stats/OtherStats";
export interface StatsViewerProps {
  result: Result;
}

export default function StatsViewer(props: StatsViewerProps) {
  const { result } = props;

  return (
    <div>
      <CollapsingSection title="Death Stats">
        <DeathStats result={result} />
      </CollapsingSection>
      <CollapsingSection title="Other Stats">
        <OtherStats result={result} />
      </CollapsingSection>
    </div>
  );
}
