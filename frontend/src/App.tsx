import { useState, useEffect } from "react";
import axios from "axios";
import { Divider, Grid, Typography } from "@mui/material";
import { Container } from "@mui/system";

import "./App.css";

function App() {
  const [thumbnails, setThumbnails] = useState<Array<string>>([]);

  useEffect(() => {
    const fetchThumbnails = async () => {
      const res = await axios.get("/media/thumbnail/all");
      for (let thumbnail of res.data.thumbnails) {
        setThumbnails((images) => [...images, "data:image/png;base64," + thumbnail]);
      }
    };

    fetchThumbnails();
  }, []);

  return (
    <Container maxWidth="xl">
      <Typography variant="h3">Photolens</Typography>
      <Divider />
      <Grid container spacing={1} className="grid">
        {thumbnails.map((image) => (
          <Grid item>
            <img src={image} key={image} />
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}

export default App;
