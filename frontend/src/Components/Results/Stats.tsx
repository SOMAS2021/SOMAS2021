import ReportCard from "./ReportCard";
import {
  Average,
  DeathsPerAgent,
  Max,
  Min,
  UtilityOnDeath,
  AverageUtilityPerAgent,
  AverageAgeUponDeath,
  ParseMessageStats,
  ParseTreatyAcceptanceStats,
} from "../../Helpers/Utils";
import { Result } from "../../Helpers/Result";
import BarChart from "./Graphs/BarChart";
import MultiBarChart from "./Graphs/MultiBarChart";
import LineChart from "./Graphs/LineChart";

interface StatsViewerProps {
  result: Result;
}

export default function StatsViewer(props: StatsViewerProps) {
  const { result } = props;
  let deaths = DeathsPerAgent(result.deaths);
  let utilityPerAgent = AverageUtilityPerAgent(result.utility);
  let utilityUponDeath = UtilityOnDeath(result.utility);
  let averageAgeUponDeath = AverageAgeUponDeath(result.deaths);
  let messageStats = ParseMessageStats(result.messages);
  let treatyAcceptanceStats = ParseTreatyAcceptanceStats(result.messages);
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
          <LineChart
            yAxis={result.food.map((f) => f.food)}
            xAxis={result.food.map((f) => f.tick)}
            graphTitle="Total Food on Platform per Tick"
          />
        </div>
        <div className="col-lg-6">
          <BarChart
            yAxis={Object.values(utilityPerAgent)}
            xAxis={Object.keys(utilityPerAgent)}
            graphTitle="Average utility per Agent type"
          />
        </div>
        <div className="col-lg-6">
          <BarChart
            yAxis={Object.values(averageAgeUponDeath)}
            xAxis={Object.keys(averageAgeUponDeath)}
            graphTitle="Average age upon death per Agent type"
          />
        </div>
        <div className="col-lg-6">
          <MultiBarChart xAxis={result.messages.atypes} data={treatyAcceptanceStats} />
        </div>
        <div className="row">
          <MultiBarChart xAxis={result.messages.atypes} data={messageStats} />
        </div>
      </div>
    </div>
  );
}
