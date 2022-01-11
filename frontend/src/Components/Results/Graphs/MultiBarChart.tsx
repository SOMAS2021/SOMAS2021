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

interface MultiBarChartProps {
  xAxis: any[];
  data: any[];
}

export default function MultiBarChart(props: MultiBarChartProps) {
  const { xAxis, data } = props;
  return (
    <Bar
      data={{
        labels: xAxis,
        datasets: data,
      }}
    />
  );
}
