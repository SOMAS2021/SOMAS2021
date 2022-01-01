import { Button, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
import { displayParams } from "../ParameterLabels";
import { SimConfig } from "../SimConfig";
import "./Settings.css";

export default function NewRun(state: any) {
  const config = state[0];
  const setConfig = state[1];

  function configHandler<Key extends keyof SimConfig>(value: number, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
  }

  function SubmitState() {
    const configJSON = JSON.stringify(config);
    SubmitSimulation(configJSON);
  }

  return (
    <div
      className="modal custom fade"
      id="exampleModal"
      data-backdrop="false"
      tabIndex={-1}
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="bp3-heading">New Run Configuration</h5>
            <Button className="bp3-minimal close" icon="cross" text="" data-dismiss="modal" aria-label="Close" />
          </div>
          <div className="modal-body">
            {displayParams.slice(0, 11).map((i) => (
              <FormGroup {...i}>
                <NumericInput placeholder={config[i.key]} onValueChange={(value) => configHandler(value, i.key)} />
              </FormGroup>
            ))}
            <Button
              className="bp3-minimal"
              icon="cog"
              text="Advanced Settings"
              data-toggle="modal"
              data-target="#testModal"
              data-dismiss="modal"
            />
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
      </div>
    </div>
  );
}
