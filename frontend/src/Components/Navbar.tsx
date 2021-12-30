import { Navbar, Alignment, Button } from "@blueprintjs/core";

export default function Nav() {
  return (
    <Navbar fixedToTop={true} className="bp3-dark">
      <Navbar.Group align={Alignment.LEFT}>
        <Navbar.Heading>Platform Dashboard</Navbar.Heading>
        <Navbar.Divider />
        <Button className="bp3-minimal" icon="home" text="Home" />
        <Button className="bp3-minimal" icon="document" text="Files" />
      </Navbar.Group>
    </Navbar>
  );
}
