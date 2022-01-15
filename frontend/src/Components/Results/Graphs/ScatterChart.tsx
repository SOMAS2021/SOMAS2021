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
import zoomPlugin from "chartjs-plugin-zoom";
import { Max, Min } from "../../../Helpers/Utils";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, zoomPlugin);

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
        plugins: {
          zoom: {
            limits: {
              y: { min: Min(yAxis), max: Max(yAxis) },
              x: { min: Min(xAxis), max: Max(xAxis) },
            },
            zoom: {
              drag: {
                enabled: true,
              },
              wheel: {
                enabled: true,
              },
              mode: "xy",
            },
          },
        },
        scales: {
          y: {
            ticks: {
              // Include a dollar sign in the ticks
              callback: function (value, index, ticks) {
                // when the floored value is the same as the value we have a whole number
                return (Math.round((value as number) * 100) / 100).toFixed(2) + " " + (yUnit ? yUnit : "");
              },
            },
          },
          x: {
            ticks: {
              // Include a dollar sign in the ticks
              callback: function (value, index, ticks) {
                return (Math.round((value as number) * 100) / 100).toFixed(2) + " " + (xUnit ? xUnit : "");
              },
            },
          },
        },
      }}
    />
  );
}
