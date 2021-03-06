import { FormGroup, H5, NumericInput, Switch } from "@blueprintjs/core";
import { useState } from "react";
import { Parameter } from "./ParameterLabels";

export default function TowerFood(props: any) {
  const { config, configHandler } = props;
  const [disableTotalFood, setDisableTotalFood] = useState(true);
  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Food Information</H5>
      <div className="row">
        <div className="col-lg-6 d-flex justify-content-center" key="Switch">
          <FormGroup>
            <Switch
              label="Use Food Per Agent"
              checked={config.UseFoodPerAgentRatio}
              onChange={(value) => {
                setDisableTotalFood((value.target as HTMLInputElement).checked);
                configHandler((value.target as HTMLInputElement).checked, "UseFoodPerAgentRatio");
              }}
            />
          </FormGroup>
        </div>
        {foodParams.map((i) =>
          i.key === "FoodOnPlatform" ? (
            <div className="col-lg-6 d-flex justify-content-center" key={i.key}>
              <FormGroup {...i} disabled={disableTotalFood}>
                <NumericInput
                  disabled={disableTotalFood}
                  placeholder={config[i.key].toString()}
                  onValueChange={(value) => configHandler(value, i.key)}
                  min={i.min}
                />
              </FormGroup>
            </div>
          ) : i.key === "FoodPerAgentRatio" ? (
            <div className="col-lg-6 d-flex justify-content-center" key={i.key}>
              <FormGroup {...i} disabled={!disableTotalFood}>
                <NumericInput
                  disabled={!disableTotalFood}
                  placeholder={config[i.key].toString()}
                  onValueChange={(value) => configHandler(value, i.key)}
                  min={i.min}
                />
              </FormGroup>
            </div>
          ) : (
            <div className="col-lg-6 d-flex justify-content-center" key={i.key}>
              <FormGroup {...i}>
                <NumericInput
                  placeholder={config[i.key].toString()}
                  onValueChange={(value) => configHandler(value, i.key)}
                  min={i.min}
                />
              </FormGroup>
            </div>
          )
        )}
      </div>
    </div>
  );
}

const foodParams: Parameter[] = [
  {
    helperText: "Food on the platform at the beginning of each day",
    label: "Initial Food On The Platform",
    labelFor: "text-input",
    labelInfo: "",
    key: "FoodOnPlatform",
    min: 1,
  },
  {
    helperText: "Maximum Food Intake For An Agent",
    label: "Max Food Intake",
    labelFor: "text-input",
    labelInfo: "",
    key: "MaxFoodIntake",
    min: 1,
  },
  {
    helperText: "Food on the platform at the beginning of each day",
    label: "Food Per Agent",
    labelFor: "text-input",
    labelInfo: "",
    key: "FoodPerAgentRatio",
    min: 1,
  },
];
