import styled from "@emotion/styled";

import logo from "../assets/logo.png";

const StyledLogo = styled.img`
  max-width: 200px;
  margin: 1rem 0;
`;

export default function Logo() {
  return <StyledLogo src={logo} alt="Photolens" />;
}
