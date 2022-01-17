import { Button, Colors } from "@blueprintjs/core";
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
import { uniqueId } from "@blueprintjs/core/lib/esm/common/utils";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, zoomPlugin);

interface ScatterChartProps {
  xAxis: number[];
  yAxis: number[];
  graphTitle: string;
  yUnit?: string;
  xUnit?: string;
  color?: string;
}

export default function ScatterChart(props: ScatterChartProps) {
  const { xAxis, yAxis, graphTitle, yUnit, xUnit } = props;
  var { color } = props;
  if (color === undefined) {
    color = Colors.BLUE1;
  }
  const uuid = uniqueId("somas");
  const resetChart = () => {
    const x = ChartJS.getChart(uuid);
    if (x !== undefined) {
      x.resetZoom();
    } else {
      console.log("undefined chart " + uuid);
    }
  };
  return (
    <div className="d-flex justify-content-center">
      <div className="border rounded" style={{ width: "95%", padding: 10 }}>
        <Scatter
          id={uuid}
          data={{
            datasets: [
              {
                data: Array.from(xAxis.keys()).map((i) => {
                  return { x: xAxis[i], y: yAxis[i] };
                }),
                showLine: true,
                label: graphTitle,
                pointBackgroundColor: color,
                backgroundColor: color,
                borderColor: color,
              },
            ],
          }}
          options={{
            plugins: {
              zoom: {
                pan: {
                  enabled: true,
                  onPanStart({ chart, point }) {
                    const area = chart.chartArea;
                    if (point.x < area.left || point.x > area.right || point.y < area.top || point.y > area.bottom) {
                      return false; // abort
                    }
                  },
                  mode: "xy",
                },
                limits: {
                  y: { min: Min(yAxis), max: Max(yAxis) },
                  x: { min: Min(xAxis), max: Max(xAxis) },
                },
                zoom: {
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
                  callback: function (value, index, ticks) {
                    return (Math.round((value as number) * 100) / 100).toFixed(2) + " " + (yUnit ? yUnit : "");
                  },
                },
              },
              x: {
                ticks: {
                  callback: function (value, index, ticks) {
                    return (Math.round((value as number) * 100) / 100).toFixed(2) + " " + (xUnit ? xUnit : "");
                  },
                },
              },
            },
          }}
        />
        <Button text="Reset Zoom" onClick={() => resetChart()} />
      </div>
    </div>
  );
}
