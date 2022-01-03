import { Button, Checkbox, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
import { params } from "../ParameterLabels";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./Settings.css";

export default function Settings(state: any) {
  const config = state[0];
  const setConfig = state[1];

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
    console.log(config);
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
            {params.map((i) => (
              <FormGroup {...i}>
                <NumericInput placeholder={config[i.key]} onValueChange={(value) => configHandler(value, i.key)} min={i.min} />
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
      </div>
    </div>
  );
}
