import { Button, Card, Elevation, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
import { displayParams } from "../ParameterLabels";
import { SimConfig } from "../SimConfig";
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
          <Card interactive={true} elevation={Elevation.TWO}>
            {displayParams.map((i) => (
              <FormGroup {...i}>
                <NumericInput placeholder={config[i.key]} onValueChange={(value) => configHandler(value, i.key)} />
              </FormGroup>
            ))}
            <Button intent="danger" className="close" icon="cross" text="Cancel" data-dismiss="modal" />
            <Button
              intent="success"
              icon="build"
              text="Submit job to backend"
              data-dismiss="modal"
              onClick={() => SubmitState()}
            />
          </Card>
        </div>
        <div className="modal-footer"></div>
      </div>
    </div>
  );
}
