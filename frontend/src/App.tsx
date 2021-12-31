import "./App.css";
import Navbar from "./Components/Navbar";
import NewRun from "./Components/NewRun";
import Results from "./Components/Results";
import Sidebar from "./Components/Sidebar";
function App() {
  var a = 3
  return (
    <div>
      <Navbar />
      <div style={{ height: 50 }}></div>
      <div className="row">
        <div className="col-lg-2">
          <Sidebar />
        </div>
        <div className="col-lg-10">
          <Results />
        </div>
      </div>
      <NewRun />
    </div>
  );
}

export default App;
