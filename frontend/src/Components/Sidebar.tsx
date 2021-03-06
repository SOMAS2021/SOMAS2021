import { H6, Menu, MenuDivider, MenuItem, Spinner } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import { GetLogs } from "../Helpers/API";

interface SideBarProps {
  activeLog: string;
  setActiveLog: React.Dispatch<React.SetStateAction<string>>;
  logsSub: boolean;
}

export default function Sidebar(props: SideBarProps) {
  const { activeLog, setActiveLog, logsSub } = props;
  const [loading, setLoading] = useState(true);
  const [logs, setLogs] = useState<string[]>([]);

  useEffect(() => {
    GetLogs(loading)
      .then((logs) => {
        setLogs(logs);
        setLoading(false);
      })
      .catch((_) => setLoading(false));
  }, [logsSub, loading]);
  return (
    <div
      style={{
        overflowY: "scroll",
        overflowX: "hidden",
        height: "95vh",
        textAlign: "left",
        padding: "10px 0px",
        borderColor: "#EBF1F5",
        borderWidth: 0,
        borderRightWidth: 10,
        borderStyle: "solid",
      }}
      className="hide-scrollbar"
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
          <H6 className="text-center" style={{ padding: 10 }}>
            Congratulations! You've reached the end of the directory
          </H6>
        </Menu>
      )}
      {loading && <Spinner intent="primary" />}
    </div>
  );
}
