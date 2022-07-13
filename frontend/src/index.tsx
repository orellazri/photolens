import "normalize.css/normalize.css";
import "@fontsource/roboto/300.css";
import "@fontsource/roboto/400.css";
import "@fontsource/roboto/500.css";
import "@fontsource/roboto/700.css";
import ReactDOM from "react-dom/client";
import axios from "axios";

import App from "./App";

const root = ReactDOM.createRoot(document.getElementById("root") as HTMLElement);

axios.defaults.baseURL = "http://localhost:5000";
axios.defaults.headers.post["Content-Type"] = "application/json";

root.render(<App />);
