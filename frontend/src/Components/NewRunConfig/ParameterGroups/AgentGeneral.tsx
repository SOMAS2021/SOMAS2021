import { H5, FormGroup, NumericInput } from "@blueprintjs/core";
import { Parameter } from "./ParameterLabels";

export default function AgentGeneral(props: any) {
  const { config, configHandler } = props;
  var classNameString = "col-lg-4 d-flex justify-content-center";

  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Agent General</H5>
      <div className="row">
        {agentParams.map((i) => (
          <div className={classNameString} key={i.key}>
            <FormGroup {...i}>
              <NumericInput
                placeholder={config[i.key].toString()}
                onValueChange={(value) => configHandler(value, i.key)}
                min={i.min}
              />
            </FormGroup>
          </div>
        ))}
      </div>
    </div>
  );
}

const agentParams: Parameter[] = [
  {
    helperText: "",
    label: "Agents Initial HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentHP",
    min: 1,
  },
  {
    helperText: "",
    label: "Agents Max HP",
    labelFor: "text-input",
    labelInfo: "",
    key: "maxHP",
    min: 1,
  },
  {
    helperText: "",
    label: "Agents Per Floor",
    labelFor: "text-input",
    labelInfo: "",
    key: "AgentsPerFloor",
    min: 1,
  },
];
