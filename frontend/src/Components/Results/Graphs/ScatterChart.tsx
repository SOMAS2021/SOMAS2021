import { Colors } from "@blueprintjs/core";
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

interface ScatterChartProps {
  xAxis: number[];
  yAxis: number[];
  graphTitle: string;
  yUnit?: string;
  xUnit?: string;
}

export default function ScatterChart(props: ScatterChartProps) {
  const { xAxis, yAxis, graphTitle, yUnit, xUnit } = props;
  return (
    <Scatter
      data={{
        datasets: [
          {
            data: Array.from(xAxis.keys()).map((i) => {
              return { x: xAxis[i], y: yAxis[i] };
            }),
            showLine: true,
            label: graphTitle,
            pointBackgroundColor: Colors.BLUE1,
            backgroundColor: Colors.BLUE1,
            borderColor: Colors.BLUE1,
          },
        ],
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
