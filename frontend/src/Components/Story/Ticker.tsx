import { H4, Button, Divider } from "@blueprintjs/core";
import { StoryFoodLog, StoryMessageLog, StoryDeathLog, StoryPlatformLog } from "../../Helpers/Logging/StoryLog";

interface TickerProps {
  tick: number;
  setTick: React.Dispatch<React.SetStateAction<number>>;
  maxTick: number;
}

export function Ticker(props: TickerProps) {
  const { tick, setTick, maxTick } = props;
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
    </div>
  );
}
