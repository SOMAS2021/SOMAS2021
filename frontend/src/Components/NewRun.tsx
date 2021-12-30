import { Button } from "@blueprintjs/core";
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
          <div className="modal-body">Form should go here</div>
          <div className="modal-footer">
            <Button intent="danger" className="close" icon="cross" text="Cancel" data-dismiss="modal" />
            <Button intent="success" icon="build" text="Submit job to backend" data-dismiss="modal" onClick={() => showToast("Job submitted successfully to backend!", "success")}/>
          </div>
        </div>
      </div>
    </div>
  );
}
