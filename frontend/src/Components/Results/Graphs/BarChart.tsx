import { Colors } from "@blueprintjs/core";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend,
} from "chart.js";
import { Bar } from "react-chartjs-2";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, BarElement, Title, Tooltip, Legend);

interface BarChartProps {
  xAxis: any[];
  yAxis: number[];
  graphTitle: string;
}

export default function BarChart(props: BarChartProps) {
  const { xAxis, yAxis, graphTitle } = props;

  return (
    <Bar
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
