import { Menu, MenuItem } from "@blueprintjs/core";

export default function Sidebar() {
  function handleClick() {}
  return (
    <div style={{ overflowY: "scroll", height: "95vh", textAlign: "left"}}>
      <Menu>
        {[...range(1, 100)].map((i) => (
          <MenuItem icon="document" onClick={handleClick} text={`This is log ${i}`} active={window.location.pathname === `/${i}`} href={`/${i}`} />
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
