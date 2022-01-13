import { Button, Collapse, Divider, Intent, Pre } from "@blueprintjs/core";
import { useState } from "react";
import {
  StoryDeathLog,
  StoryFoodLog,
  StoryLog,
  StoryMessageLog,
  StoryPlatformLog,
} from "../../Helpers/Logging/StoryLog";
import { Ticker } from "./Ticker";

interface StoryViewerProps {
  story: StoryLog[];
}

export default function StoryViewer(props: StoryViewerProps) {
  const { story } = props;
  const [isOpen, setIsOpen] = useState(false);
  return (
    <div style={{ margin: "10px 0px" }}>
      <Button
        intent={isOpen ? Intent.PRIMARY : Intent.WARNING}
        onClick={() => setIsOpen(!isOpen)}
        style={{ width: 200 }}
      >
        {isOpen ? "Hide" : "Show"} Story
      </Button>
      <Collapse isOpen={isOpen} keepChildrenMounted={true}>
        <Pre>
          <StoryController story={story} />
        </Pre>
      </Collapse>
    </div>
  );
}

interface StoryControllerProps {
  story: StoryLog[];
}

function StoryController(props: StoryControllerProps) {
  const { story } = props;
  const [tick, setTick] = useState(1);
  const maxTick = Math.max(...story.map((e) => e.tick));

  return (
    <div style={{ overflow: "hidden" }}>
      <Ticker tick={tick} setTick={setTick} maxTick={maxTick} />
      <Divider />
      <div style={{ height: "45vh", overflowY: "scroll" }}>
        {story.map((log, index) => {
          if (log.tick === tick) {
            switch (log.msg) {
              case "food":
                const f = log as StoryFoodLog;
                return (
                  <div key={index}>
                    Agent {f.atype} took {f.foodTaken} food on floor {f.floor} and left {f.foodLeft}
                  </div>
                );
              case "message":
                const m = log as StoryMessageLog;
                return (
                  <div key={index}>
                    Agent {m.atype} on Floor {m.floor} sent a message {m.mtype} targeting Floor {m.target} with{" "}
                    {m.mcontent === "" && "no "}content {m.mcontent}
                  </div>
                );
              case "death":
                const d = log as StoryDeathLog;
                return (
                  <div key={index}>
                    Agent {d.atype} died at age {d.age} on floor {d.floor}
                  </div>
                );
              case "platform":
                const p = log as StoryPlatformLog;
                return (
                  <div key={index}>
                    The platform moved from floor {p.floor - 1} to floor {p.floor}
                  </div>
                );
              default:
                return <div key={index}>oups {log.msg}</div>;
            }
          }
          return null;
        })}
      </div>
    </div>
  );
}
