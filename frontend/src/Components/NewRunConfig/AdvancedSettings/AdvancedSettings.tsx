import { Button, Card, Elevation, FormGroup, NumericInput } from "@blueprintjs/core";
import { SubmitSimulation } from "../NewRunState";
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
            <FormGroup helperText="In Days..." label="Simulation Length" labelFor="text-input">
              <NumericInput
                placeholder={config["SimDays"]}
                onValueChange={(value) => configHandler(value, "SimDays")}
              />
            </FormGroup>
            <FormGroup
              helperText="Agents can do one/two actions per 'Tick'"
              label="'Ticks' Per Floor"
              labelFor="text-input"
            >
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "TicksPerFloor")} />
            </FormGroup>
            <FormGroup
              helperText="Food on the platfrom at the beginnning of each day"
              label="Initial Food On The Platform"
              labelFor="text-input"
            >
              <NumericInput placeholder="100" onValueChange={(value) => configHandler(value, "FoodOnPlatform")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 1 in the Tower"
              label="Agents of Team 1"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team1Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 2 in the Tower"
              label="Agents of Team 2"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team2Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 3 in the Tower"
              label="Agents of Team 3"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team3Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 4 in the Tower"
              label="Agents of Team 4"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team4Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 5 in the Tower"
              label="Agents of Team 5"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team5Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 6 in the Tower"
              label="Agents of Team 6"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team6Agents")} />
            </FormGroup>
            <FormGroup
              helperText="Number of Agents from Team 7 in the Tower"
              label="Agents of Team 7"
              labelFor="text-input"
            >
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "Team7Agents")} />
            </FormGroup>
            <FormGroup helperText="Number of Random Agents in the Tower" label="Random Agents" labelFor="text-input">
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "RandomAgents")} />
            </FormGroup>
            <FormGroup helperText="" label="Agents Max HP" labelFor="text-input">
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "AgentHP")} />
            </FormGroup>
            <FormGroup helperText="" label="Agents Per Floor" labelFor="text-input">
              <NumericInput placeholder="2" onValueChange={(value) => configHandler(value, "AgentsPerFloor")} />
            </FormGroup>
            <FormGroup helperText="In Days..." label="Reshuffle Period" labelFor="text-input">
              <NumericInput placeholder="2" onValueChange={(value) => configHandler(value, "ReshuffleDays")} />
            </FormGroup>
            <FormGroup helperText="" label="Weak Level" labelFor="text-input">
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "weakLevel")} />
            </FormGroup>
            <FormGroup helperText="" label="Width" labelFor="text-input">
              <NumericInput placeholder="45" onValueChange={(value) => configHandler(value, "width")} />
            </FormGroup>
            <FormGroup helperText="" label="Tau" labelFor="text-input">
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "tau")} />
            </FormGroup>
            <FormGroup helperText="" label="HP Required from Critical to Weak" labelFor="text-input">
              <NumericInput placeholder="2" onValueChange={(value) => configHandler(value, "hpReqCToW")} />
            </FormGroup>
            <FormGroup helperText="Upper bound" label="Critical HP" labelFor="text-input">
              <NumericInput placeholder="5" onValueChange={(value) => configHandler(value, "hpCritical")} />
            </FormGroup>
            <FormGroup helperText="" label="Days at Critical" labelFor="text-input">
              <NumericInput placeholder="3" onValueChange={(value) => configHandler(value, "maxDayCritical")} />
            </FormGroup>
            <FormGroup helperText="" label="HP Loss Base" labelFor="text-input">
              <NumericInput placeholder="10" onValueChange={(value) => configHandler(value, "HPLossBase")} />
            </FormGroup>
            <FormGroup helperText="" label="HP Loss Slope" labelFor="text-input">
              <NumericInput placeholder="0.25" onValueChange={(value) => configHandler(value, "HPLossSlope")} />
            </FormGroup>
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
