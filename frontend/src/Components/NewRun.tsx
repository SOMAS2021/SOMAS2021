import { Button, FormGroup, InputGroup } from "@blueprintjs/core";
import { showToast } from "./Toaster";

export default function NewRun() {
  return (
    <div
      className="modal custom fade"
      id="exampleModal"
      data-backdrop="false"
      tabIndex={-1}
      aria-labelledby="staticBackdropLabel"
      aria-hidden="true"
    >
      <div className="modal-dialog">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="bp3-heading">New Run Configuration</h5>
            <Button className="bp3-minimal close" icon="cross" text="" data-dismiss="modal" aria-label="Close" />
          </div>
          <div className="modal-body">
            <FormGroup
              helperText="Some helper test"
              label="Parameters"
              labelFor="text-input"
              labelInfo="(required)"
            >
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
              <InputGroup id="text-input" placeholder="Placeholder text" style ={{margin: "10px 0px"}}/>
            </FormGroup>
          </div>
          <div className="modal-footer">
            <Button intent="danger" className="close" icon="cross" text="Cancel" data-dismiss="modal" />
            <Button
              intent="success"
              icon="build"
              text="Submit job to backend"
              data-dismiss="modal"
              onClick={() => showToast("Job submitted successfully to backend!", "success")}
            />
          </div>
        </div>
      </div>
    </div>
  );
}
