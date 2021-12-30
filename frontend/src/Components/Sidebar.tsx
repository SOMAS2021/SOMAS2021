import { Menu, MenuItem } from "@blueprintjs/core";
import { useState } from "react";

export default function Sidebar() {
  const [log, setLog] = useState(0);
  return (
    <div
      style={{
        overflowY: "scroll",
        overflowX: "hidden",
        height: "95vh",
        textAlign: "left",
        padding: "10px 0px",
      }}
    >
      <Menu>
        {[...range(1, 100)].map((i) => (
          <MenuItem
            icon="document"
            onClick={() => {
              setLog(i);
            }}
            text={`This is log ${i}`}
            active={log === i}
          />
        ))}
      </Menu>
    </div>
  );
}

function range(start: number, end: number) {
  return Array(end - start + 1)
    .fill(0)
    .map((_, idx) => start + idx);
}
