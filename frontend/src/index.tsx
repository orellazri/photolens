import "normalize.css/normalize.css";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { BrowserRouter, Link, Route, Routes } from "react-router-dom";
import { Container } from "@mui/system";
import { Button, Divider, Stack } from "@mui/material";
import ReactDOM from "react-dom/client";
import axios from "axios";

import "./index.css";
import Logo from "./components/Logo";
import HomePage from "./pages/HomePage";
import SettingsPage from "./pages/SettingsPage";

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

// Set global variables
global.API_URL = "http://localhost:5000";

// Set axios defaults
axios.defaults.baseURL = global.API_URL;
axios.defaults.headers.post["Content-Type"] = "application/json";

root.render(
  <BrowserRouter>
    {/* Container */}
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
  </BrowserRouter>
);
