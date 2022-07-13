import { useState, useEffect } from "react";
import axios from "axios";
import { Divider, Grid, Typography } from "@mui/material";
import { Container } from "@mui/system";

import "./App.css";

type Thumbnail = {
  id: number;
  thumbnail: string;
};

function App() {
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);

  useEffect(() => {
    const fetchThumbnails = async () => {
      const res = await axios.get("/media/thumbnail/all");
      for (let item of res.data.data) {
        setThumbnails((images) => [...images, { id: item.id, thumbnail: "data:image/png;base64," + item.thumbnail }]);
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
          <Grid item key={image.id}>
            <img src={image.thumbnail} />
          </Grid>
        ))}
      </Grid>
    </Container>
  );
}

export default App;
