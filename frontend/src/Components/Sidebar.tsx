import { Menu, MenuDivider, MenuItem, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetLogs } from "../Helpers/API";

interface SideBarProps {
  activeLog: string;
  setActiveLog: React.Dispatch<React.SetStateAction<string>>;
}

export default function Sidebar(props: SideBarProps) {
  const { activeLog, setActiveLog } = props;
  const [loading, setLoading] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    setLoading(true);
    GetLogs()
      .then((logs) => {
        setLogs(logs);
        setLoading(false);
      })
      .catch((_) => setLoading(false));
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
