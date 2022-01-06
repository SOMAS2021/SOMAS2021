import { showToast } from "../Components/Toaster";
import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { GetStoryLogs, StoryLog } from "./Logging/StoryLog";
import { Result } from "./Result";
import { GetSimConfig, SimConfig } from "./SimConfig";

function endpoint(req: string) {
  return (process.env.REACT_APP_DEV ? "http://localhost:9000/" : window.location) + req;
}

function parseResponse(res: any, key: string) {
  return res[key].map(function (e: any) {
    return JSON.parse(e);
  });
}

export function GetLogs(): Promise<string[]> {
  return new Promise<string[]>((resolve, reject) => {
    showToast("Loading logs in progress", "primary");
    fetch(endpoint("directory"))
      .then(async (res) => {
        if (res.status !== 200) {
          showToast(`Loading logs failed. (${res.status}) ${await res.text()}`, "danger", 5000);
          reject(res);
        }
        res
          .json()
          .then((res) => {
            showToast("Loading logs completed", "success");
            resolve(res["FolderNames"]);
          })
          .catch((err) => {
            showToast(`Loading logs: failed. ${err}`, "danger", 5000);
            reject(err);
          });
      })
      .catch((err) => {
        showToast(`Loading logs: failed. ${err}`, "danger", 5000);
        reject(err);
      });
  });
}

export function GetResult(filename: string): Promise<Result> {
  return new Promise<Result>(async (resolve, reject) => {
    // promises
    var promises: Promise<any>[] = [];
    // deaths
    var deaths: DeathLog[] = [];
    promises.push(GetDeathLogs(filename).then((d) => (deaths = d)));

    // foods
    var foods: FoodLog[] = [];
    promises.push(GetFoodLogs(filename).then((f) => (foods = f)));

    // config
    var config: SimConfig = undefined!;
    promises.push(GetSimConfig(filename).then((c) => (config = c)));

    // story
    var story: StoryLog[] = [];
    promises.push(GetStoryLogs(filename).then((s) => (story = s)));

    // all
    Promise.all(promises).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
        food: foods,
        config: config,
        story: story,
      })
    );
  });
}

export function GetFile(filename: string, logtype: string): Promise<any> {
  return new Promise<any>((resolve, reject) => {
    var requestOptions = {
      method: "POST",
      body: JSON.stringify({ LogFileName: filename, LogType: logtype }),
    };

    fetch(endpoint("read"), requestOptions)
      .then((response) => response.json())
      .catch((error) => {
        showToast(`Loading file: failed. ${error}`, "danger", 5000);
        reject(error);
      })
      .then((result) => resolve(parseResponse(result, "Log")))
      .catch((error) => {
        showToast(`Loading file: failed. ${error}`, "danger", 5000);
        reject(error);
      });
  });
}

export function Simulate(config: SimConfig): Promise<boolean> {
  return new Promise<boolean>((resolve, reject) => {
    const requestOptions = {
      method: "POST",
      body: JSON.stringify(config),
    };
    showToast("Job submitted successfully to backend!", "success");
    fetch(endpoint("simulate"), requestOptions)
      .then(function (response) {
        response.json().then((res) => console.log(res));
        resolve(true);
      })
      .catch(function (error) {
        console.log("There has been a problem with submitting the simulation: " + error.message);
        reject(error);
      });
  });
}
