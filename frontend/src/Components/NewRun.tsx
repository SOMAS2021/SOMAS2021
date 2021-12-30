import { Button, FormGroup, NumericInput } from "@blueprintjs/core";
import InitiateSimConfig, { simConfig } from "./Simconfig";
import { showToast } from "./Toaster";


function request(configJSON: string){
  const requestOptions = {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: configJSON
  };
  var host = window.location.protocol + "//" + window.location.host;
  fetch(`${host}/simulate`, requestOptions)
}

export default function NewRun() {

  const [config, setConfig] = InitiateSimConfig()
  
  const configHandler = <Key extends keyof simConfig>(value: number, keyString: any) => {
    var key:Key = keyString // converting keyString to type Key
    var c = config;
    c[key]=value;
    setConfig(c)
    console.log(c)
    
  }
  const submitSimulation = () => {
    const configJSON = JSON.stringify(config);
    console.log(configJSON)
    request(configJSON)
    showToast("Job submitted successfully to backend!", "success")
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
          <FormGroup
              helperText="Length of Simulation in days..."
              label="Simulation Length"
              labelFor="text-input"
              labelInfo="(required)"
          >
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "SimDays")} />
          </FormGroup>
          </div>
          <div className="modal-footer">
            <Button intent="danger" className="close" icon="cross" text="Cancel" data-dismiss="modal" />
            <Button
              intent="success"
              icon="build"
              text="Submit job to backend"
              data-dismiss="modal"
              onClick={() => submitSimulation()}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
