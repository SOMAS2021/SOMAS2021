import { Button, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
import { advancedParams, params } from "../ParameterLabels";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./AdvancedSettings.css";

export default function AdvancedSettingsMenu(state: any) {
  const config = state[0];
  const setConfig = state[1];

  function configHandler<Key extends keyof SimConfig>(value: number, keyString: any) {
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
            {params.map((i) => (
              <FormGroup {...i}>
                <NumericInput placeholder={config[i.key]} onValueChange={(value) => configHandler(value, i.key)} />
              </FormGroup>
            ))}
            {advancedParams.map((i) => (
              <FormGroup {...i}>
                <NumericInput placeholder={config[i.key]} onValueChange={(value) => configHandler(value, i.key)} />
              </FormGroup>
            ))}
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
