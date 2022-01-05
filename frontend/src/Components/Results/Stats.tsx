import ReportCard from "./ReportCard";
import { Average, Max, Min } from "../../Helpers/Utils";
import { Result } from "../../Helpers/Result";

interface StatsViewerProps {
  result: Result;
}

export default function StatsViewer(props: StatsViewerProps) {
  const { result } = props;
  return (
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
          title={Max(result.deaths.map((d) => d.ageUponDeath)).toString()}
        />
      </div>
      <div className="col-lg-6">
        <ReportCard
          description="Min agent age upon death"
          title={Min(result.deaths.map((d) => d.ageUponDeath)).toString()}
        />
      </div>
    </div>
  );
}
