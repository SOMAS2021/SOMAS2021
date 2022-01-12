import { H3, Button, Collapse } from "@blueprintjs/core";
import { useState } from "react";

interface CollapsingSectionProps {
  children: JSX.Element;
  title: string;
}

export default function CollapsingSection(props: CollapsingSectionProps) {
  const { children, title } = props;
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div>
      <H3 style={{ paddingTop: 20 }}>
        <Button
          icon={isOpen ? "chevron-down" : "chevron-right"}
          intent="none"
          onClick={() => setIsOpen(!isOpen)}
          style={{ marginRight: 10 }}
        ></Button>
        {title}
      </H3>
      <div style={{ margin: "10px 0px" }}>
        <Collapse isOpen={isOpen}>{children}</Collapse>
      </div>
    </div>
  );
}
