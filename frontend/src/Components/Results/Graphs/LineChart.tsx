import { Card, Elevation, H1, H6 } from "@blueprintjs/core";
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import {Line} from 'react-chartjs-2';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface LineChartProps {
  xAxis: any[]
  yAxis: number[]
  graphTitle: string
}

const data = {
  // labels: ["Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul"],
  
  datasets: [
    {
      label: "My First dataset",
      data: [1500000, 3900000, 3000000, 4100000, 2300000, 1800000, 2000000]
    }
  ]
};
export default function LineChart(props: LineChartProps) {
  const { xAxis, yAxis , graphTitle} = props;
  return (
    <Line data={{
            labels: xAxis,
            datasets: [
                {
                    label: graphTitle,
                    data: yAxis,
                    backgroundColor: 'blue',
                }
            ]
        }}
    />
  );
}

