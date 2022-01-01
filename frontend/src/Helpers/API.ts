import { showToast } from "../Components/Toaster";

function endpoint(req: string) {
  return "http://localhost:9000/" + req;
  return window.location + req;
}

export function GetLogs(): Promise<string[]> {
  console.log(endpoint("directory"))
  return new Promise<string[]>((resolve, reject) => {
    showToast("Loading logs in progress", "primary");
    fetch(endpoint("directory"))
      .then(async (res) => {
        if (res.status !== 200) {
          showToast(`Loading logs failed. (${res.status}) ${await res.text()}`, "danger", 5000);
          reject(res)
        }
        res.json().then((res) => {
          showToast("Loading logs completed", "success");
          resolve(res["FolderNames"])
        });
      })
      .catch((err) => {
        showToast(`Loading logs: failed. ${err}`, "danger", 5000);
        reject(err)
      });
  });
}
