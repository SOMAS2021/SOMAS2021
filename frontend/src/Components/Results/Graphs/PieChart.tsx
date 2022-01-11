// import { Colors } from "@blueprintjs/core";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  ArcElement,
} from "chart.js";
import { Pie } from "react-chartjs-2";

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement, ArcElement, Title, Tooltip, Legend);

interface PieChartProps {
  labels: any[];
  data: any[];
}

export default function PieChart(props: PieChartProps) {
  const { labels, data } = props;

  return (
    <Pie
      data={{
        labels: labels,
        datasets: data,
      }}
    />
  );
}
