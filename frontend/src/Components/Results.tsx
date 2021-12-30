import { Button, Card, Elevation } from "@blueprintjs/core";

export default function Results() {
  return (
    <div style={{ padding: 20 }}>
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
    </div>
  );
}

function CardExample() {
  return (
    <Card
      interactive={true}
      elevation={Elevation.TWO}
      style={{ marginTop: 20 }}
    >
      <h5>
        <a href="#">100%</a>
      </h5>
      <p>Description</p>
    </Card>
  );
}
