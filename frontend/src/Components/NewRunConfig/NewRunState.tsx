import { showToast } from "../Toaster";
import AdvancedSettingsMenu from "./AdvancedSettings/AdvancedSettings";
import NewRun from "./Settings/Settings";
import InitConfigState from "../../Helpers/SimConfig";

function request(configJSON: string) {
  const requestOptions = {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: configJSON,
  };
  var host = window.location.protocol + "//" + window.location.host;
  fetch(`${host}/simulate`, requestOptions);
}

export function SubmitSimulation(configJSON: string) {
  console.log(configJSON);
  request(configJSON);
  showToast("Job submitted successfully to backend!", "success");
}

export default function NewRunState() {
  // config state declaration
  const [config, setConfig] = InitConfigState();

  return (
    <div>
      <AdvancedSettingsMenu {...[config, setConfig]} />
      <NewRun {...[config, setConfig]} />
    </div>
  );
}