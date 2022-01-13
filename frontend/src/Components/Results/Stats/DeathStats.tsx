import { DeathLog } from "../../../Helpers/Logging/Death";
import { DeathsPerAgent, UtilityOnDeath, Average, Max, Min, AverageAgeUponDeath } from "../../../Helpers/Utils";
import BarChart from "../Graphs/BarChart";
import MultiScatterChart from "../Graphs/MultiScatterChart";
import ScatterChart from "../Graphs/ScatterChart";
import ReportCard from "../ReportCard";
import { StatsViewerProps } from "../Stats";

export default function DeathStats(props: StatsViewerProps) {
  const { result } = props;
  let deaths = DeathsPerAgent(result.deaths);
  let utilityUponDeath = UtilityOnDeath(result.utility);
  let averageAgeUponDeath = AverageAgeUponDeath(result.deaths);

  // Multi death graph data
  var yAxisDeathAgents = [];
  var xAxisDeathAgents = [];
  var titleDeathAgents = ["Team1Agent1", "Team2", "Team3", "Team4", "Team5", "Team6", "Team7", "RandomAgent"];
  var color = ["#1F4B99", "#447C9F", "#7CAAA2", "#CCD3A1", "#F6C880", "#DE944D", "#C06126", "#9E2B0E"];

  for (let i = 0; i < titleDeathAgents.length; i++) {
    const agentType = titleDeathAgents[i];
    const checkType = (d: DeathLog) => {
      return d.agentType === agentType;
    };
    yAxisDeathAgents.push([0].concat(result.deaths.filter(checkType).map((d) => d.cumulativeDeaths)));
    xAxisDeathAgents.push([0].concat(result.deaths.filter(checkType).map((d) => d.day)));
    console.log(result.deaths.filter(checkType));
  }

  console.log(xAxisDeathAgents, yAxisDeathAgents);

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
        <div className="col-lg-12">
          <ScatterChart
            yAxis={[0].concat(
              result.deaths
                .map((d) => d.cumulativeDeaths)
                .concat(result.deaths.length > 0 ? result.deaths[result.deaths.length - 1].cumulativeDeaths : 0)
            )}
            xAxis={[0].concat(result.deaths.map((d) => d.day).concat(result.config.SimDays))}
            graphTitle="Cumulative deaths per day"
          />
        </div>
        <div className="col-lg-12">
          <MultiScatterChart
            xAxis={xAxisDeathAgents}
            yAxis={yAxisDeathAgents}
            graphTitle={titleDeathAgents}
            color={color}
          />
        </div>
      </div>
    </div>
  );
}
