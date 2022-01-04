import { Button, Collapse, H2, H4, H5, Pre, Intent, Divider } from "@blueprintjs/core";
import blob1 from "../../assets/blobs/blob1.png";
import blob2 from "../../assets/blobs/blob2.png";
import blob3 from "../../assets/blobs/blob3.png";
import blob4 from "../../assets/blobs/blob4.png";
import blob5 from "../../assets/blobs/blob5.png";
import blob6 from "../../assets/blobs/blob6.png";
import blob7 from "../../assets/blobs/blob7.png";
import blob8 from "../../assets/blobs/blob8.png";
import { useState } from "react";
import { SimConfig } from "../../Helpers/SimConfig";

interface ConfigInfoProps {
  config: SimConfig;
}

export default function ConfigInfo(props: ConfigInfoProps) {
  const { config } = props;

  const [isOpen, setIsOpen] = useState(false);

  return (
    <div style={{ margin: "10px 0px" }}>
      <Button intent={isOpen ? Intent.PRIMARY : Intent.WARNING} onClick={() => setIsOpen(!isOpen)}>
        {isOpen ? "Hide" : "Show"} config
      </Button>
      <Collapse isOpen={isOpen}>
        <Pre>
          <Params config={config} />
          <Divider />
          <Agents config={config} />
        </Pre>
      </Collapse>
    </div>
  );
}

function Params(props: ConfigInfoProps) {
  const { config } = props;

  return (
    <div className="row" style={{ padding: "10px 20px" }}>
      <H4>Parameters</H4>
      {mainParams(config).map((param, index) => (
        <div className="col-md-2" key={index}>
          <H2>{param.value}</H2>
          <H5>{param.name}</H5>
        </div>
      ))}
      <div className="col-md-2">
        <H2>{config.UseFoodPerAgentRatio ? config.FoodPerAgentRatio : config.FoodOnPlatform}</H2>
        <H5>{config.UseFoodPerAgentRatio ? "Food per agent" : "Food on platform"}</H5>
      </div>
    </div>
  );
}

function Agents(props: ConfigInfoProps) {
  const { config } = props;

  return (
    <div className="row" style={{ padding: "10px 20px" }}>
      <H4>Agents</H4>
      {blobs(config).map((blob, index) => (
        <div className="col-md-3" key={index}>
          <H2>
            <img
              src={blob.blob}
              style={{ paddingRight: 10, opacity: blob.count === 0 ? 0.5 : 1 }}
              alt={`Blob image for ${blob.name}`}
            ></img>
            x{blob.count}
          </H2>
          <H5>({blob.name})</H5>
        </div>
      ))}
    </div>
  );
}

function mainParams(config: SimConfig) {
  return [
    {
      name: "Agent HP",
      value: config.AgentHP,
    },
    {
      name: "Max HP",
      value: config.maxHP,
    },
    {
      name: "Agents / Floor",
      value: config.AgentsPerFloor,
    },
    {
      name: "Ticks / Floor",
      value: config.TicksPerFloor,
    },
    {
      name: "Days",
      value: config.SimDays,
    },
  ];
}

function blobs(config: SimConfig) {
  return [
    {
      count: config.Team1Agents,
      name: "Team 1",
      blob: blob1,
    },
    {
      count: config.Team2Agents,
      name: "Team 2",
      blob: blob2,
    },
    {
      count: config.Team3Agents,
      name: "Team 3",
      blob: blob3,
    },
    {
      count: config.Team4Agents,
      name: "Team 4",
      blob: blob4,
    },
    {
      count: config.Team5Agents,
      name: "Team 5",
      blob: blob5,
    },
    {
      count: config.Team6Agents,
      name: "Team 6",
      blob: blob6,
    },
    {
      count: config.Team7Agents,
      name: "Team 7",
      blob: blob7,
    },
    {
      count: config.RandomAgents,
      name: "Random Agents",
      blob: blob8,
    },
  ];
}
