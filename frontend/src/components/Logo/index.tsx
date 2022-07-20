import { Link } from "react-router-dom";

import "./style.css";
import logo from "../../assets/logo.png";

export default function Logo() {
  return (
    <Link to="/">
      <img src={logo} alt="Photolens" className="logo" />
    </Link>
  );
}
