import { Divider, H5, Pre } from "@blueprintjs/core";
import { useEffect, useState } from "react";
import {
  GetStoryLogs,
  StoryDeathLog,
  StoryFoodLog,
  StoryLog,
  StoryMessageLog,
  StoryPlatformLog,
} from "../../Helpers/Logging/StoryLog";
import { Result } from "../../Helpers/Result";
import { Ticker } from "./Ticker";

interface StoryViewerProps {
  result: Result;
}

export default function StoryViewer(props: StoryViewerProps) {
  const { result } = props;
  return (
    <div style={{ margin: "10px 0px" }}>
      {result.simStatus.maxStoryTick > -1 ? (
        <Pre>
          <StoryController title={result.title} maxTick={result.simStatus.maxStoryTick} />
        </Pre>
      ) : result.config.LogStory ? (
        <H5>Simulation still in progress</H5>
      ) : (
        <H5>
          <i>Story unavailable</i>
        </H5>
      )}
    </div>
  );
}

interface StoryControllerProps {
  title: string;
  maxTick: number;
}

function StoryController(props: StoryControllerProps) {
  const { title, maxTick } = props;
  const [tick, setTick] = useState(1);
  const [story, setStory] = useState<StoryLog[]>([]);

  useEffect(() => {
    GetStoryLogs(title, tick).then((story) => setStory(story));
  }, [tick, title]);

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
