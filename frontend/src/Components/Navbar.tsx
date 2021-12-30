import { Navbar, Alignment, Button } from "@blueprintjs/core";
import logo from "../assets/experiment.png";
export default function Nav() {
  return (
    <Navbar fixedToTop={true} className="bp3-dark">
      <Navbar.Group align={Alignment.LEFT}>
        <Navbar.Heading><img src={logo} alt="logo" height={25} style={{paddingRight: 10}}/>Platform Dashboard</Navbar.Heading>
        <Navbar.Divider />
        <Button className="bp3-minimal" icon="search-around" text="New Run" />
      </Navbar.Group>
    </Navbar>
  );
}
