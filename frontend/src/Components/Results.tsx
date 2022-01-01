import { Card, Elevation, H3, H5, H6 } from "@blueprintjs/core";

interface ResultsProps {
  logName: string;
}

export default function Results(props: ResultsProps) {
  const { logName } = props;
  return (
    <div style={{ padding: 20 }}>
      {logName !== "" ? (
        <>
          <H3>{logName}</H3>
          <Cards />
        </>
      ) : (
        <H6>
          <i>Select an exsiting simulation result to view results</i>
        </H6>
      )}
    </div>
  );
}

function Cards() {
  return (
    <div className="row">
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-2">
        <CardExample />
      </div>
      <div className="col-lg-6">
        <CardExample />
      </div>
      <div className="col-lg-4">
        <CardExample />
      </div>
    </div>
  );
}

function CardExample() {
  return (
    <Card interactive={true} elevation={Elevation.TWO} style={{ marginTop: 20 }}>
      <H5 style={{ color: "#1F4B99" }}>100%</H5>
      <p>Description</p>
    </Card>
  );
}
