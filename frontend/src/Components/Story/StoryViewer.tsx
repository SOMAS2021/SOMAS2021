import { Divider, H2, H5, Icon, Pre } from "@blueprintjs/core";
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
import GetBlob from "./StoryHelpers";
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

  // todo: max floor count
  const floors = 4;

  // Tower
  const tower: JSX.Element[] = [];
  for (let f = 0; f < floors; f++) {
    tower.push(
      <div
        key={f}
        style={{ position: "absolute", top: f * 100, border: "solid 3px black", width: "80%", height: "90px" }}
      >
        <H2 style={{ paddingLeft: 20 }}>Floor {f + 1}</H2>
      </div>
    );
  }

  return (
    <div style={{ overflow: "hidden" }}>
      <Ticker tick={tick} setTick={setTick} maxTick={maxTick} />
      <Divider />
      {/* Draw story */}
      <div style={{ height: "45vh", overflowY: "scroll", position: "relative", top: 0, bottom: 0 }}>
        {story.map((log, index) => {
          if (log.tick === tick) {
            switch (log.msg) {
              case "food":
                const f = log as StoryFoodLog;
                return (
                  <></>
                  // <div key={index}>
                  //   Agent {f.atype} took {f.foodTaken} food on floor {f.floor} and left {f.foodLeft}
                  // </div>
                );
              case "message":
                const m = log as StoryMessageLog;
                const floorDiff = m.target - m.floor;
                const neg = floorDiff < 0;
                return (
                  <>
                    <H5 style={{ position: "absolute", top: (m.floor - 1) * 100 + 15, left: 200 }}>
                      <img alt="" src={GetBlob(m.atype)} key={index} />
                      HP={m.hp}
                      <div
                        style={{
                          position: "absolute",
                          left: index * 50 + 400 + 7,
                          top: neg ? 75 * floorDiff : 20 * floorDiff,
                          width: 2,
                          border: "solid 1px black",
                          height: 110 * floorDiff * (neg ? -1 : 1),
                        }}
                      />
                      <Icon
                        icon={neg ? "arrow-up" : "arrow-down"}
                        style={{ top: 100 * floorDiff + 20, position: "absolute", left: index * 50 + 400 }}
                      />
                      <p style={{ top: 100 * floorDiff + 20, position: "absolute", left: index * 50 + 420 }}>{m.mtype}</p>
                    </H5>
                  </>
                  // <div key={index}>
                  //   Agent {m.atype} on Floor {m.floor} sent a message {m.mtype} targeting Floor {m.target} with{" "}
                  //   {m.mcontent === "" && "no "}content {m.mcontent}
                  // </div>
                );
              case "death":
                const d = log as StoryDeathLog;
                return (
                  <></>
                  // <img alt="" src={GetBlob(d.atype)} key={index} style={{position: "absolute", top: d.floor*100, left:100}} />
                  // <div key={index}>
                  //   Agent {d.atype} died at age {d.age} on floor {d.floor}
                  // </div>
                );
              case "platform":
                const p = log as StoryPlatformLog;
                return (
                  <></>
                  // <div key={index}>
                  //   The platform moved from floor {p.floor - 1} to floor {p.floor}
                  // </div>
                );
              default:
                return <div key={index}>oups {log.msg}</div>;
            }
          }
          return null;
        })}
        {tower}
      </div>
    </div>
  );
}
