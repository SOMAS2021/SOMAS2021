import { showToast } from "../Components/Toaster";
import { SimConfig } from "./SimConfig";

function endpoint(req: string) {
  return (process.env.REACT_APP_DEV ? "http://localhost:9000/" : window.location) + req;
}

function parseResponse(res: any, key: string) {
  return res[key].map(function (e: any) {
    return JSON.parse(e);
  });
}

export function GetLogs(toast: boolean = true): Promise<string[]> {
  return new Promise<string[]>((resolve, reject) => {
    toast && showToast("Loading logs in progress", "primary");
    fetch(endpoint("directory"))
      .then(async (res) => {
        if (res.status !== 200) {
          toast && showToast(`Loading logs failed. (${res.status}) ${await res.text()}`, "danger", 5000);
          reject(res);
        }
        res
          .json()
          .then((res) => {
            toast && showToast("Loading logs completed", "success");
            resolve(res["FolderNames"]);
          })
          .catch((err) => {
            toast && showToast(`Loading logs: failed. ${err}`, "danger", 5000);
            reject(err);
          });
      })
      .catch((err) => {
        toast && showToast(`Loading logs: failed. ${err}`, "danger", 5000);
        reject(err);
      });
  });
}

export function GetFile(filename: string, logtype: string, tick: number = -1): Promise<any> {
  return new Promise<any>((resolve, reject) => {
    var requestOptions = {
      method: "POST",
      body: JSON.stringify({
        LogFileName: filename,
        LogType: logtype,
        TickFilter: tick > -1,
        Tick: tick,
      }),
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
