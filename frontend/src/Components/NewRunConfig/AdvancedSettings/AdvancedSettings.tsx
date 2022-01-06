import { Button, Checkbox, FormGroup, NumericInput, Switch } from "@blueprintjs/core";
import { advancedParams, params } from "../ParameterLabels";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./AdvancedSettings.css";
import { useState } from "react";
import { Simulate } from "../../../Helpers/API";

interface AdvancedSettingsProps {
  config: SimConfig;
  setConfig: React.Dispatch<React.SetStateAction<SimConfig>>;
}

export default function AdvancedSettings(props: AdvancedSettingsProps) {
  const { config, setConfig } = props;

  const [disableTotalFood, setDisableTotalFood] = useState(Boolean);

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
  }

  return (
    <div
      className="modal custom fade"
      id="testModal"
      data-backdrop="false"
      tabIndex={-1}
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">Advanced Settings</div>
          <div className="modal-body">
            <FormGroup>
              <Switch
                label="Use Food Per Agent"
                onChange={(value) => {
                  setDisableTotalFood((value.target as HTMLInputElement).checked);
                  configHandler((value.target as HTMLInputElement).checked, "UseFoodPerAgentRatio");
                }}
              />
            </FormGroup>
            {params.map((i) =>
              i.key === "FoodOnPlatform" ? (
                <FormGroup {...i} disabled={disableTotalFood}>
                  <NumericInput
                    disabled={disableTotalFood}
                    placeholder={config[i.key].toString()}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              ) : i.key === "FoodPerAgentRatio" ? (
                <FormGroup {...i} disabled={!disableTotalFood}>
                  <NumericInput
                    disabled={!disableTotalFood}
                    placeholder={config[i.key].toString()}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              ) : (
                <FormGroup {...i}>
                  <NumericInput
                    placeholder={config[i.key].toString()}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              )
            )}
            {advancedParams.map((i) => (
              <FormGroup {...i}>
                <NumericInput
                  placeholder={config[i.key].toString()}
                  onValueChange={(value) => configHandler(value, i.key)}
                  min={i.min}
                />
              </FormGroup>
            ))}
            <FormGroup>
              <Checkbox
                label="Save Main"
                type="checkbox"
                onChange={(value) => configHandler((value.target as HTMLInputElement).checked, "LogMain")}
              />
            </FormGroup>
            <FormGroup label="File Name" labelFor="text-input" key="FileName">
              <input
                type="text"
                onChange={(value) => configHandler(value.target.value, "LogFileName")}
                placeholder="Simulation Run #"
              />
            </FormGroup>
          </div>
          <div className="modal-footer">
            <Button intent="danger" className="close" icon="cross" text="Cancel" data-dismiss="modal" />
            <Button
              intent="success"
              icon="build"
              text="Submit job to backend"
              data-dismiss="modal"
              onClick={() => Simulate(config)}
            />
          </div>
        </div>
        <div className="modal-footer"></div>
      </div>
    </div>
  );
}
