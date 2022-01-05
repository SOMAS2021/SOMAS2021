import { Button, Collapse, Divider, H4, Intent, Pre } from "@blueprintjs/core";
import { useState } from "react";
import {
  StoryDeathLog,
  StoryFoodLog,
  StoryLog,
  StoryMessageLog,
  StoryPlatformLog,
} from "../../Helpers/Logging/StoryLog";

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
    <div>
      <div className="row">
        <div className="col-lg-6">
          <H4>
            Tick {tick} / {maxTick}
          </H4>
        </div>
        <div className="col-lg-6" style={{ textAlign: "right" }}>
          <Button icon="arrow-left" onClick={() => setTick(Math.max(1, tick - 100))} disabled={tick <= 1} />
          <Button icon="double-chevron-left" onClick={() => setTick(Math.max(1, tick - 10))} disabled={tick <= 1} />
          <Button icon="chevron-left" onClick={() => setTick(Math.max(1, tick - 1))} disabled={tick <= 1} />
          <Button
            icon="chevron-right"
            onClick={() => setTick(Math.min(maxTick, tick + 1))}
            disabled={tick >= maxTick}
          />
          <Button
            icon="double-chevron-right"
            onClick={() => setTick(Math.min(maxTick, tick + 10))}
            disabled={tick >= maxTick}
          />
          <Button
            icon="arrow-right"
            onClick={() => setTick(Math.min(maxTick, tick + 100))}
            disabled={tick >= maxTick}
          />
        </div>
      </div>
      <Divider />
      <div>
        {story.map((log, index) => {
          if (log.tick === tick) {
            switch (log.msg) {
              case "food":
                const f = log as StoryFoodLog;
                return <div key={index}>{f.foodTaken}</div>;
              case "message":
                const m = log as StoryMessageLog;
                return <div key={index}>{m.mtype}</div>;
              case "death":
                const d = log as StoryDeathLog;
                return <div key={index}>{d.atype}</div>;
              case "platform":
                const p = log as StoryPlatformLog;
                return <div key={index}>{p.floor}</div>;
              default:
                return <div key={index}>oups</div>;
            }
          }
          return null;
        })}
      </div>
    </div>
  );
}
