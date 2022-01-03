import { Button, Checkbox, FormGroup, NumericInput, Switch } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
import { advancedParams, params } from "../ParameterLabels";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./AdvancedSettings.css";
import { useState } from "react";

export default function AdvancedSettingsMenu(state: any) {
  const config = state[0];
  const setConfig = state[1];

  const [disableTotalFood, setDisableTotalFood] = useState(Boolean);

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
  }

  function SubmitState() {
    const configJSON = JSON.stringify(state[0]);
    SubmitSimulation(configJSON);
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
                    placeholder={config[i.key]}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              ) : i.key === "FoodPerAgentRatio" ? (
                <FormGroup {...i} disabled={!disableTotalFood}>
                  <NumericInput
                    disabled={!disableTotalFood}
                    placeholder={config[i.key]}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              ) : (
                <FormGroup {...i}>
                  <NumericInput
                    placeholder={config[i.key]}
                    onValueChange={(value) => configHandler(value, i.key)}
                    min={i.min}
                  />
                </FormGroup>
              )
            )}
            {advancedParams.map((i) => (
              <FormGroup {...i}>
                <NumericInput
                  placeholder={config[i.key]}
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
              onClick={() => SubmitState()}
            />
          </div>
        </div>
        <div className="modal-footer"></div>
      </div>
    </div>
  );
}
