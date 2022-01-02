import { showToast } from "../Components/Toaster";
import { DeathLog, GetDeathLogs } from "./Logging/Death";
import { FoodLog, GetFoodLogs } from "./Logging/Food";
import { Result } from "./Result";

function endpoint(req: string) {
  return (process.env.DEV ? "http://localhost:9000/" : window.location) + req;
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

    // all
    Promise.all(promises).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
        food: foods,
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
