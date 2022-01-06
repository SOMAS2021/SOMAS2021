import { Button, H3 } from "@blueprintjs/core";
import { SimConfig } from "../../../Helpers/SimConfig";
import "./AdvancedSettings.css";
import { Simulate } from "../../../Helpers/API";
import TowerFood from "../ParameterGroups/TowerFood";
import TowerLength from "../ParameterGroups/TowerLength";
import AgentTypesParams from "../ParameterGroups/AgentTypes";
import AgentGeneral from "../ParameterGroups/AgentGeneral";
import FileName from "../ParameterGroups/FileName";

interface AdvancedSettingsProps {
  config: SimConfig;
  setConfig: React.Dispatch<React.SetStateAction<SimConfig>>;
}

export default function AdvancedSettings(props: AdvancedSettingsProps) {
  const { config, setConfig } = props;

  function configHandler<Key extends keyof SimConfig>(value: any, keyString: any) {
    var key: Key = keyString; // converting keyString to type Key
    config[key] = value;
    setConfig(config);
  }

  return (
    <div
      className="modal custom fade"
      id="advancedSettingsModal"
      data-backdrop="false"
      tabIndex={-1}
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <H3 className="text-center">Advanced Settings</H3>
          </div>
          <div className="modal-body">
            <TowerFood config={config} configHandler={configHandler} />
            <TowerLength config={config} configHandler={configHandler} advanced={true} />
            <AgentTypesParams config={config} configHandler={configHandler} />
            <AgentGeneral config={config} configHandler={configHandler} advanced={true} />
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
