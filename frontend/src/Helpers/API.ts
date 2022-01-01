import { showToast } from "../Components/Toaster";
import { DeathLog } from "./Logs/Death";
import { Result } from "./Result";

const DEV = true;

function endpoint(req: string) {
  return (DEV ? "http://localhost:9000/" : window.location) + req;
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
        res.json().then((res) => {
          showToast("Loading logs completed", "success");
          resolve(res["FolderNames"]);
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
    var deaths: DeathLog[] = [];
    const getDeaths = GetDeathLogs(filename).then((d) => (deaths = d));

    Promise.all([getDeaths]).then((_) =>
      resolve({
        title: filename,
        deaths: deaths,
      })
    );
  });
}

function GetFile(filename: string, logtype: string): Promise<any> {
  return new Promise<any>((resolve, reject) => {
    var requestOptions = {
      method: "POST",
      body: JSON.stringify({ LogFileName: filename, LogType: logtype }),
    };

    fetch("http://localhost:9000/read", requestOptions)
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

function GetDeathLogs(filename: string): Promise<DeathLog[]> {
  return new Promise<DeathLog[]>((resolve, reject) => {
    GetFile(filename, "death")
      .then((deaths) =>
        resolve(
          deaths.map(function (e: any) {
            const d: DeathLog = {
              cumulativeDeaths: e["cumulativeDeaths"],
              agentType: e["agent_type"],
              tick: e["tick"],
              day: e["day"],
            };
            return d;
          })
        )
      )
      .catch((err) => reject(err));
  });
}
