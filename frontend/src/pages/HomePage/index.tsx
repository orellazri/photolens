import { Container } from "@mui/system";

import "./style.css";
import Logo from "../../components/Logo";
import Gallery from "../../components/Gallery";

export default function HomePage() {
  return (
    <Container maxWidth="xl">
      <Logo />
      <Gallery />
    </Container>
  );
}
