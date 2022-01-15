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
import { Min, Max } from "../../../Helpers/Utils";
import zoomPlugin from "chartjs-plugin-zoom";
import { uniqueId } from "@blueprintjs/core/lib/esm/common/utils";
import { Button } from "@blueprintjs/core";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend, zoomPlugin);

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
      <div className="border rounded" style={{ width: "90%" }}>
        <Scatter
          id={uuid}
          data={{
            datasets: datasets,
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
                  y: { min: Min(yAxis.map((i) => Min(i))), max: Max(yAxis.map((i) => Max(i))) },
                  x: { min: Min(xAxis.map((i) => Min(i))), max: Max(xAxis.map((i) => Max(i))) },
                },
                zoom: {
                  wheel: {
                    enabled: true,
                  },
                  mode: "x",
                },
              },
            },
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
        <Button text="Reset Zoom" onClick={() => resetChart()} />
      </div>
    </div>
  );
}
