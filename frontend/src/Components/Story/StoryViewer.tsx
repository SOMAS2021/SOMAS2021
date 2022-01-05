import { Button } from "@blueprintjs/core";
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
  const [tick, setTick] = useState(1);
  return (
    <div>
      <Button icon="arrow-left" onClick={() => setTick(tick - 1)}/>
      <Button icon="arrow-right" onClick={() => setTick(tick + 1)}/>
      Tick {tick}
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
