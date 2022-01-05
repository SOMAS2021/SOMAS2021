import { FormGroup, H5, NumericInput, Switch } from "@blueprintjs/core";
import { Parameter } from "./ParameterLabels";

export default function TowerLength(props: any) {
  const { config, configHandler, advanced } = props;
  var classNameString = "col-lg-4 d-flex justify-content-center";
  if (advanced === true) {
    classNameString = "col-lg-6 d-flex justify-content-center";
  }
  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Tower Information</H5>
      <div className="row">
        {lengthParams.map((i) => (
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
      <div className="row">
        {lengthParamsAdv.map((i) =>
          advanced === true ? (
            <div className={classNameString} key={i.key}>
              <FormGroup {...i}>
                <NumericInput
                  placeholder={config[i.key].toString()}
                  onValueChange={(value) => configHandler(value, i.key)}
                  min={i.min}
                />
              </FormGroup>
            </div>
          ) : (
            <div key={i.key}></div>
          )
        )}
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
];

export const lengthParamsAdv: Parameter[] = [
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
