import { FormGroup, H5, NumericInput } from "@blueprintjs/core";
import { Parameter } from "./ParameterLabels";

export default function TowerLength(props: any) {
  const { config, configHandler } = props;
  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Tower Information</H5>
      <div className="row">
        {lengthParams.map((i) => (
          <div className="col-lg-6 d-flex justify-content-center" key={i.key}>
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

export const lengthParams: Parameter[] = [
  {
    helperText: "In Days...",
    label: "Simulation Length",
    labelFor: "text-input",
    labelInfo: "",
    key: "SimDays",
    min: 1,
  },
  {
    helperText: "In Days...",
    label: "Reshuffle Period",
    labelFor: "text-input",
    labelInfo: "",
    key: "ReshuffleDays",
    min: 1,
  },
  {
    helperText: "Agents can do one/two actions per 'Tick'",
    label: "'Ticks' Per Floor",
    labelFor: "text-input",
    labelInfo: "",
    key: "TicksPerFloor",
    min: 1,
  },
  {
    helperText: "Base this on the size of your simulation.",
    label: "TimeOut",
    labelFor: "text-input",
    labelInfo: "",
    key: "SimTimeoutSeconds",
    min: 1,
  },
];
