import "normalize.css/normalize.css";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import ReactDOM from "react-dom/client";
import axios from "axios";

import "./index.css";
import HomePage from "./pages/HomePage";

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

// Set global variables
global.API_URL = "http://localhost:5000";

// Set axios defaults
axios.defaults.baseURL = global.API_URL;
axios.defaults.headers.post["Content-Type"] = "application/json";

root.render(
  <BrowserRouter>
    <Routes>
      <Route path="/" element={<HomePage />} />
    </Routes>
  </BrowserRouter>
);
