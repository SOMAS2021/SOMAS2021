import { Button } from "@blueprintjs/core";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./Settings.css";
import { useState } from "react";
import { Simulate } from "../../../Helpers/API";
import TowerFood from "../ParameterGroups/TowerFood";
import TowerLength from "../ParameterGroups/TowerLength";
import AgentGeneral from "../ParameterGroups/AgentGeneral";
import AgentTypesParams from "../ParameterGroups/AgentTypes";
import LogDescription from "../ParameterGroups/LogDescription";

interface SettingsProps {
  config: SimConfig;
  setConfig: React.Dispatch<React.SetStateAction<SimConfig>>;
}

export default function Settings(props: SettingsProps) {
  const { config, setConfig } = props;

  const [disableTotalFood, setDisableTotalFood] = useState(Boolean);

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
    console.log(config);
  }

  return (
    <div
      className="modal custom fade"
      id="settingsModal"
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
            <TowerFood config={config} configHandler={configHandler} />
            <TowerLength config={config} configHandler={configHandler} advanced={false} />
            <AgentTypesParams config={config} configHandler={configHandler} />
            <AgentGeneral config={config} configHandler={configHandler} advanced={false} />
            <LogDescription configHandler={configHandler} />
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
      </div>
    </div>
  );
}
