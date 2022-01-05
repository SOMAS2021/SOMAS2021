import AdvancedSettingsMenu from "./AdvancedSettings/AdvancedSettings";
import Settings from "./Settings/Settings";
import InitConfigState from "../../Helpers/SimConfig";

export default function NewRunState() {
  // config state declaration
  const [config, setConfig] = InitConfigState();

  return (
    <div>
      <AdvancedSettingsMenu config={config} setConfig={setConfig} />
      <Settings config={config} setConfig={setConfig} />
    </div>
  );
}
