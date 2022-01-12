import {
  DeathsPerAgent,
  AverageUtilityPerAgent,
  ParseMessageStats,
  ParseTreatyAcceptanceStats,
  Average,
} from "../../../Helpers/Utils";
import BarChart from "../Graphs/BarChart";
import LineChart from "../Graphs/LineChart";
import MultiBarChart from "../Graphs/MultiBarChart";
import ReportCard from "../ReportCard";
import { StatsViewerProps } from "../Stats";

export default function OtherStats(props: StatsViewerProps) {
  const { result } = props;
  let deaths = DeathsPerAgent(result.deaths);
  let utilityPerAgent = AverageUtilityPerAgent(result.utility);
  let messageStats = ParseMessageStats(result.messages);
  let treatyAcceptanceStats = ParseTreatyAcceptanceStats(result.messages);
  return (
    <div className="row">
      <div className="col-lg-6">
        <ReportCard
          description="Average food on platform per tick"
          title={Average(result.food.map((f) => f.food))
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
          <MultiBarChart xAxis={result.messages.atypes} data={treatyAcceptanceStats} />
        </div>
        <div className="row">
          <MultiBarChart xAxis={result.messages.atypes} data={messageStats} />
        </div>
      </div>
    </div>
  );
}
