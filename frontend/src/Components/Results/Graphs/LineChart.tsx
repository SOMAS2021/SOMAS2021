import { Card, Elevation, H1, H6, Intent, Colors } from "@blueprintjs/core";
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
import { Line } from "react-chartjs-2";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, Title, Tooltip, Legend);

interface LineChartProps {
  xAxis: any[];
  yAxis: number[];
  graphTitle: string;
}

export default function LineChart(props: LineChartProps) {
  const { xAxis, yAxis, graphTitle } = props;
  return (
    <Line
      data={{
        labels: xAxis,
        datasets: [
          {
            label: graphTitle,
            data: yAxis,
            backgroundColor: Colors.BLUE1,
          },
        ],
      }}
    />
  );
}
