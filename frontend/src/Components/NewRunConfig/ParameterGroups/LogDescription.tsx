import { H5, FormGroup, Switch } from "@blueprintjs/core";

export default function LogDescription(props: any) {
  const { configHandler } = props;

  return (
    <div style={{ paddingTop: 20 }}>
      <H5 className="text-center">Log Description</H5>
      <div className="row">
        <div className="col-lg-6 d-flex justify-content-center">
          <FormGroup>
            <Switch
              label="Save Main"
              onChange={(value) => {
                configHandler((value.target as HTMLInputElement).checked, "LogMain");
              }}
            />
          </FormGroup>
        </div>
        <div className="col-lg-6 d-flex justify-content-center">
          <FormGroup label="File Name" labelFor="text-input" key="FileName">
            <input
              type="text"
              onChange={(value) => configHandler(value.target.value, "LogFileName")}
              placeholder="Simulation Run #"
            />
          </FormGroup>
        </div>
      </div>
    </div>
  );
}
