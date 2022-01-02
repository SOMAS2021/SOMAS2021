import { showToast } from "../Toaster";
import AdvancedSettingsMenu from "./AdvancedSettings/AdvancedSettings";
import Settings from "./Settings/Settings";
import InitConfigState from "../../Helpers/SimConfig";

function request(configJSON: string) {
  const requestOptions = {
    method: "POST",
    body: configJSON,
  };
  var host = window.location.protocol + "//" + window.location.host;
  fetch(`${host}/simulate`, requestOptions)
    .then(function (response) {
      return response.json();
    })
    .catch(function (error) {
      console.log("There has been a problem with your fetch operation: " + error.message);
    });
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
      <Settings {...[config, setConfig]} />
    </div>
  );
}
