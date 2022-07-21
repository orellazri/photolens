import { useState } from "react";
import { Alert, Button, Snackbar } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";

import "./style.css";

export default function SettingsPage() {
  const [error, setError] = useState<string>("");
  const [success, setSuccess] = useState<string>("");

  const [isProcessing, setIsProcessing] = useState<boolean>(false);

  const handleProcessMedia = () => {
    try {
      setIsProcessing(true);
      axios.get("/media/process").then(() => {
        setSuccess("Finished processing media");
        setIsProcessing(false);
      });
    } catch (e) {
      console.error("Could not process media! " + e);
      setError("Could not process media");
    }
  };

  return (
    <Box>
      <Button variant="contained" onClick={handleProcessMedia} disabled={isProcessing}>
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
