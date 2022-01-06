import ReportCard from "./ReportCard";
import { Average, DeathsPerAgent } from "../../Helpers/Utils";
import { Result } from "../../Helpers/Result";
import BarChart from "./Graphs/BarChart";
import LineChart from "./Graphs/LineChart";

interface StatsViewerProps {
  result: Result;
}

export default function StatsViewer(props: StatsViewerProps) {
  const { result } = props;
  let deaths = DeathsPerAgent(result.deaths);
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
          title={result.deaths.length === 0 ? "0" : Math.max(...result.deaths.map((d) => d.ageUponDeath)).toString()}
        />
      </div>
      <div className="col-lg-6">
        <ReportCard
          description="Min agent age upon death"
          title={result.deaths.length === 0 ? "0" : Math.min(...result.deaths.map((d) => d.ageUponDeath)).toString()}
        />
      </div>
      <div className="row">
        <div className="col-lg-6">
          <BarChart yAxis={Object.values(deaths)} xAxis={Object.keys(deaths)} graphTitle="Deaths per Agent type" />
        </div>
        <div className="col-lg-6">
          <LineChart
            yAxis={result.food.map((f) => f.food)}
            xAxis={result.food.map((f) => f.tick)}
            graphTitle="Total Food on Platform per Tick"
          />
        </div>
      </div>
    </div>
  );
}
