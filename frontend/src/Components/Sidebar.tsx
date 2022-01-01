import { Menu, MenuDivider, MenuItem, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { showToast } from "./Toaster";

export default function Sidebar() {
  const [activeLog, setActiveLog] = useState("");
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    showToast("Loading logs in progress", "primary");
    setLoading(true);
    fetch("http://localhost:9000/directory")
      .then(async (res) => {
        if(res.status != 200){
          showToast(`Loading logs failed. (${res.status}) ${await res.text()}`, "danger", 5000);
          setLoading(false);
        }
        res.json().then((res) => {
          setLogs(res["FolderNames"]);
          showToast("Loading logs completed", "success");
          setLoading(false);
        })
      })
      .catch((err) => {
        showToast(`Loading logs: failed. ${err}`, "danger", 5000);
        setLoading(false)
      });
  }, []);
  return (
    <div
      style={{
        overflowY: "scroll",
        overflowX: "hidden",
        height: "95vh",
        textAlign: "left",
        padding: "10px 0px",
        // backgroundColor: "#EBF1F5",
      }}
    >
      {!loading && (
        <Menu>
          {logs.map((log, index) => (
            <div key={index}>
              <MenuItem
                icon="document"
                onClick={() => {
                  setActiveLog(log);
                }}
                text={`${log}`}
                active={activeLog === log}
              />
              <MenuDivider />
            </div>
          ))}
        </Menu>
      )}
      {loading && <Spinner intent="primary" />}
    </div>
  );
}
