import Settings from "./Settings/Settings";
import InitConfigState from "../../Helpers/SimConfig";

interface NewRunStateProps {
  onSubmit: () => void;
}

export default function NewRunState(props: NewRunStateProps) {
  const { onSubmit } = props;
  // config state declaration
  const [config, setConfig] = InitConfigState();

  return (
    <div>
      <Settings config={config} setConfig={setConfig} onSubmit={onSubmit} />
    </div>
  );
}
