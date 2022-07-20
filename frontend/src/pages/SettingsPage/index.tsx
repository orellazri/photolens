import { useState } from "react";
import { Alert, Button, Snackbar } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";

import "./style.css";

export default function SettingsPage() {
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<string>("");

  const handleProcessMedia = async () => {
    try {
      await axios.get("/media/process");
      setSuccess("Finished processing media");
    } catch (e) {
      console.error("Could not process media! " + e);
      setError("Could not process media");
    }
  };

  return (
    <Box>
      <Button variant="contained" onClick={handleProcessMedia}>
        Process media
      </Button>

      {/* Error snackbar */}
      <Snackbar
        open={error !== ""}
        onClose={() => {
          setError("");
        }}
        autoHideDuration={5000}
      >
        <Alert severity="error">{error}</Alert>
      </Snackbar>

      {/* Success snackbar */}
      <Snackbar
        open={success !== ""}
        onClose={() => {
          setSuccess("");
        }}
        autoHideDuration={5000}
      >
        <Alert severity="success">{success}</Alert>
      </Snackbar>
    </Box>
  );
}
