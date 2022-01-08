import { H5, FormGroup, NumericInput, Switch } from "@blueprintjs/core";
import { Parameter } from "./ParameterLabels";

export default function AgentTypesParams(props: any) {
  const { config, configHandler } = props;
  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Agent Composition</H5>
      <div className="row">
        {agentTypeParams.map((i) => (
          <div className="col-lg-4 d-flex justify-content-center" key={i.key}>
            <FormGroup {...i}>
              <NumericInput
                placeholder={config[i.key].toString()}
                onValueChange={(value) => configHandler(value, i.key)}
                min={i.min}
              />
            </FormGroup>
          </div>
        ))}
        <div className="col-lg-4 d-flex justify-content-center" key="RandomReplacementAgents">
          <FormGroup style={{ padding: 20 }}>
            <Switch
              label="Random Replacement Agents"
              onChange={(value) => {
                configHandler((value.target as HTMLInputElement).checked, "RandomReplacementAgents");
              }}
            />
          </FormGroup>
        </div>
      </div>
    </div>
  );
}

const agentTypeParams: Parameter[] = [
  {
    helperText: "Number of Team 1 Agents in the Tower",
    label: "Agents of Team 1",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team1Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 2 Agents in the Tower",
    label: "Agents of Team 2",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team2Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 3 Agents in the Tower",
    label: "Agents of Team 3",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team3Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 4 Agents in the Tower",
    label: "Agents of Team 4",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team4Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 5 Agents in the Tower",
    label: "Agents of Team 5",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team5Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 6 Agents in the Tower",
    label: "Agents of Team 6",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team6Agents",
    min: 0,
  },
  {
    helperText: "Number of Team 7 Agents in the Tower",
    label: "Agents of Team 7",
    labelFor: "text-input",
    labelInfo: "",
    key: "Team7Agents",
    min: 0,
  },
  {
    helperText: "Number of Random Agents in the Tower",
    label: "Random Agents",
    labelFor: "text-input",
    labelInfo: "",
    key: "RandomAgents",
    min: 0,
  },
];
