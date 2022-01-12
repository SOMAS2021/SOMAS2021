import { DeathsPerAgent, UtilityOnDeath, Average, Max, Min, AverageAgeUponDeath } from "../../../Helpers/Utils";
import BarChart from "../Graphs/BarChart";
import ScatterChart from "../Graphs/ScatterChart";
import ReportCard from "../ReportCard";
import { StatsViewerProps } from "../Stats";

export default function DeathStats(props: StatsViewerProps) {
  const { result } = props;
  let deaths = DeathsPerAgent(result.deaths);
  let utilityUponDeath = UtilityOnDeath(result.utility);
  let averageAgeUponDeath = AverageAgeUponDeath(result.deaths);
  return (
    <div className="row">
      <div className="col-lg-6">
        <ReportCard description="Total deaths" title={result.deaths.length.toString()} />
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
      <div className="col-lg-6">
        <ReportCard
          description="Average agent utility upon death"
          title={Average(utilityUponDeath.map((u) => u.utility))
            .toFixed(3)
            .toString()}
        />
      </div>
      <div className="row">
        <div className="col-lg-6">
          <BarChart yAxis={Object.values(deaths)} xAxis={Object.keys(deaths)} graphTitle="Deaths per Agent type" />
        </div>
        <div className="col-lg-6">
          <BarChart
            yAxis={Object.values(averageAgeUponDeath)}
            xAxis={Object.keys(averageAgeUponDeath)}
            graphTitle="Average age upon death per Agent type"
          />
        </div>
        <div className="col-lg-6">
          <BarChart yAxis={Object.values(deaths)} xAxis={Object.keys(deaths)} graphTitle="Deaths per Agent type" />
        </div>
        <div className="col-lg-6">
          <ScatterChart
            yAxis={[0].concat(result.deaths.map((d) => d.cumulativeDeaths))}
            xAxis={[0].concat(result.deaths.map((d) => d.tick))}
            graphTitle="Cumulative deaths per tick"
          />
        </div>
      </div>
    </div>
  );
}
