import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Scatter } from "react-chartjs-2";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

interface MultiScatterChartProps {
  xAxis: number[][];
  yAxis: number[][];
  graphTitle: string[];
  color: string[];
  yUnit?: string;
  xUnit?: string;
}

export default function MultiScatterChart(props: MultiScatterChartProps) {
  const { xAxis, yAxis, graphTitle, color, xUnit, yUnit } = props;
  var datasets = [];
  for (let d = 0; d < graphTitle.length; d++) {
    datasets.push({
      data: Array.from(xAxis[d].keys()).map((i) => {
        return { x: xAxis[d][i], y: yAxis[d][i] };
      }),
      showLine: true,
      label: graphTitle[d],
      pointBackgroundColor: color[d],
      backgroundColor: color[d],
      borderColor: color[d],
    });
  }
  return (
    <Scatter
      data={{
        datasets: datasets,
      }}
      options={{
        scales: {
          y: {
            ticks: {
              // Include a dollar sign in the ticks
              callback: function (value, index, ticks) {
                return value + " " + (yUnit ? yUnit : "");
              },
            },
          },
          x: {
            ticks: {
              // Include a dollar sign in the ticks
              callback: function (value, index, ticks) {
                return value + " " + (xUnit ? xUnit : "");
              },
            },
          },
        },
      }}
    />
  );
}
