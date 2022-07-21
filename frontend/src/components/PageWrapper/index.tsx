import { Link, Route, Routes } from "react-router-dom";
import { Container } from "@mui/system";
import { Button, Divider, Stack } from "@mui/material";

import "./style.css";
import Logo from "../Logo";
import HomePage from "../../pages/HomePage";
import SettingsPage from "../../pages/SettingsPage";

export default function PageWrapper() {
  return (
    <Container maxWidth="xl">
      {/* Header */}
      <Stack className="header" spacing={2}>
        <Logo />
        {/* Links */}
        <Stack direction="row" divider={<Divider orientation="vertical" flexItem />} spacing={2}>
          <Link to="/settings">
            <Button variant="outlined">Settings</Button>
          </Link>
        </Stack>
      </Stack>

      {/* Page */}
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/settings" element={<SettingsPage />} />
      </Routes>
    </Container>
  );
}
