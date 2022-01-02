import { useState } from "react";
import "./App.css";
import Navbar from "./Components/Navbar";
import Results from "./Components/Results";
import Sidebar from "./Components/Sidebar";
import NewRunState from "./Components/NewRunConfig/NewRunState";
function App() {
  const [activeLog, setActiveLog] = useState<string>("");
  return (
    <div>
      <Navbar />
      <div style={{ height: 50 }}></div>
      <div className="row">
        <div className="col-lg-2">
          <Sidebar activeLog={activeLog} setActiveLog={setActiveLog} />
        </div>
        <div className="col-lg-10">
          <Results logName={activeLog} />
        </div>
      </div>
      <NewRunState />
    </div>
  );
}

export default App;
