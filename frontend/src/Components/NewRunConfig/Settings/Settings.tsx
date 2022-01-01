import { Button, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
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
            <FormGroup helperText="In Days..." label="Simulation Length" labelFor="text-input" labelInfo="(required)">
              <NumericInput
                placeholder={config["SimDays"]}
                onValueChange={(value) => configHandler(value, "SimDays")}
              />
            </FormGroup>
            <FormGroup
              helperText="Agents can do one/two actions per 'Tick'"
              label="'Ticks' Per Floor"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "TicksPerFloor")} />
            </FormGroup>
            <FormGroup
              helperText="Food on the platfrom at the beginnning of each day"
              label="Initial Food On The Platform"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="100" onValueChange={(value) => configHandler(value, "FoodOnPlatform")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 1 in the Tower"
              label="Agents of Team 1"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team1Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 2 in the Tower"
              label="Agents of Team 2"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team2Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 3 in the Tower"
              label="Agents of Team 3"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team3Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 4 in the Tower"
              label="Agents of Team 4"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team4Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 5 in the Tower"
              label="Agents of Team 5"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team5Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 6 in the Tower"
              label="Agents of Team 6"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team6Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 7 in the Tower"
              label="Agents of Team 7"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team7Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Random Agents in the Tower"
              label="Random Agents"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "RandomAgents")} />
            </FormGroup>
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
