import { useState, useEffect } from "react";
import axios from "axios";
import { Card, CardContent, CardMedia, Divider, Grid, Skeleton, Typography } from "@mui/material";
import { Container } from "@mui/system";

import "./App.css";
import moment from "moment";

type Thumbnail = {
  id: number;
  image: string;
  createdAt: string;
};

type Metadata = {
  id: number;
  createdAt: string;
};

function App() {
  const [metadata, setMetadata] = useState<Array<Metadata>>([]);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);

  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        const res = await axios.get("/media/meta");
        setMetadata(res.data.data);
      } catch (e) {
        console.error("Could not fetch metadata! " + e);
      }
    };

    const fetchThumbnails = async () => {
      await fetchMetadata();

      await new Promise((r) => setTimeout(r, 2000));
      try {
        const res = await axios.get("/media/thumbnail/all");
        for (let item of res.data.data) {
          setThumbnails((thumbnails) => [
            ...thumbnails,
            {
              id: item.id,
              image: "data:image/png;base64," + item.thumbnail,
              createdAt: moment(item.created_at).local().format("DD/MM/YYYY HH:mm:ss"),
            },
          ]);
        }
      } catch (e) {
        console.error("Could not fetch thumbnails! " + e);
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
            thumbnails.map((thumbnail, i) => (
              <Grid item key={i}>
                <a href={`${global.API_URL}/media/${thumbnail.id}`}>
                  <Card>
                    <CardMedia component="img" height="128" image={thumbnail.image} alt={thumbnail.id.toString()} />
                    <CardContent>
                      <Typography sx={{ fontSize: 14 }} color="text.secondary" align="center" gutterBottom>
                        {thumbnail.createdAt}
                      </Typography>
                    </CardContent>
                  </Card>
                </a>
              </Grid>
            ))
          : // Show placeholder grid
            metadata.map((x, i) => (
              <Grid item key={i}>
                <Skeleton variant="rectangular" width={128} height={128} />
              </Grid>
            ))}
      </Grid>
    </Container>
  );
}

export default App;
