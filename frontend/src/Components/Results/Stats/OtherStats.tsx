import { AverageUtilityPerAgent, ParseMessageStats, ParseTreatyAcceptanceStats, Average } from "../../../Helpers/Utils";
import BarChart from "../Graphs/BarChart";
import MultiBarChart from "../Graphs/MultiBarChart";
import ScatterChart from "../Graphs/ScatterChart";
import ReportCard from "../ReportCard";
import { StatsViewerProps } from "../Stats";

export default function OtherStats(props: StatsViewerProps) {
  const { result } = props;
  let utilityPerAgent = AverageUtilityPerAgent(result.utility);
  let [msgLabels, msgValues] = ParseMessageStats(result);
  let [treatyAcceptanceLabels, treatyAcceptanceValues] = ParseTreatyAcceptanceStats(result);
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
          <ScatterChart
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
          <MultiBarChart xAxis={treatyAcceptanceLabels} data={treatyAcceptanceValues} />
        </div>
        <div className="row">
          <MultiBarChart xAxis={msgLabels} data={msgValues} />
        </div>
      </div>
    </div>
  );
}
