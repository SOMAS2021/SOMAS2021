import { Card, Elevation, H1, H5 } from "@blueprintjs/core";

interface ReportCardProps {
  title: string;
  description: string;
}

export default  function ReportCard(props: ReportCardProps) {
  const { title, description } = props;
  return (
    <Card interactive={true} elevation={Elevation.TWO} style={{ marginTop: 20 }}>
      <H1 style={{ color: "#1F4B99" }}>{title}</H1>
      <H5>{description}</H5>
    </Card>
  );
}
