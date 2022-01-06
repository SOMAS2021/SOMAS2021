import { Button, Divider, H3 } from "@blueprintjs/core";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./Settings.css";
import { Simulate } from "../../../Helpers/API";
import TowerFood from "../ParameterGroups/TowerFood";
import TowerLength from "../ParameterGroups/TowerLength";
import AgentTypesParams from "../ParameterGroups/AgentTypes";
import AgentGeneral from "../ParameterGroups/AgentGeneral";
import FileName from "../ParameterGroups/FileName";

interface SettingsInterface {
  config: SimConfig;
  setConfig: React.Dispatch<React.SetStateAction<SimConfig>>;
}

export default function Settings(props: SettingsInterface) {
  const { config, setConfig } = props;

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
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
            <H3 className="bp3-heading config-title">New Run Configuration</H3>
            <Button className="bp3-minimal close" icon="cross" text="" data-dismiss="modal" aria-label="Close" />
          </div>
          <div className="modal-body">
            <TowerFood config={config} configHandler={configHandler} />
            <Divider/>
            <TowerLength config={config} configHandler={configHandler} advanced={true} />
            <Divider/>
            <AgentTypesParams config={config} configHandler={configHandler} />
            <Divider/>
            <AgentGeneral config={config} configHandler={configHandler} advanced={true} />
            <Divider/>
            <FileName configHandler={configHandler} />
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
