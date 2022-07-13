import { useState, useEffect } from "react";
import axios from "axios";
import { Divider, Grid, Skeleton, Typography } from "@mui/material";
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
      await new Promise((r) => setTimeout(r, 2000));
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
        {thumbnails.length
          ? //  Show thumbnails grid
            thumbnails.map((image, i) => (
              <Grid item key={i}>
                <a href={`${global.API_URL}/media/${image.id}`}>
                  <img src={image.thumbnail} />
                </a>
              </Grid>
            ))
          : // Show placeholder grid
            [...Array(5)].map((x, i) => (
              <Grid item>
                <Skeleton variant="rectangular" width={128} height={128} />
              </Grid>
            ))}
      </Grid>
    </Container>
  );
}

export default App;
