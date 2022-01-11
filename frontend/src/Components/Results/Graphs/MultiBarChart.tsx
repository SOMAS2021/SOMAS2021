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
// Y axis number of messages
// X axis agent type
// let data contain the values (y-axis) and titles (message type)
// i.e. [{label1: [], data1: []}, ...]
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
