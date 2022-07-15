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

  // Fetch metadata on page load
  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        const {
          data: { data },
        } = await axios.get("/media/meta");
        for (const result of data) {
          setMetadata((metadata) => [
            ...metadata,
            {
              id: result.id,
              createdAt: result.created_at,
            },
          ]);
        }
      } catch (e) {
        console.error("Could not fetch metadata! " + e);
      }
    };

    fetchMetadata();
  }, []);

  // Fetch thumbnails after fetching
  useEffect(() => {
    const fetchThumbnails = async () => {
      try {
        for (let item of metadata) {
          const {
            data: { data },
          } = await axios.get("/media/thumbnail/" + item.id);

          setThumbnails((thumbnails) => [
            ...thumbnails,
            {
              id: data.id,
              image: "data:image/png;base64," + data.thumbnail,
              createdAt: moment(data.created_at).local().format("DD/MM/YYYY HH:mm:ss"),
            },
          ]);
        }
      } catch (e) {
        console.error("Could not fetch thumbnails! " + e);
      }
    };

    fetchThumbnails();
  }, [metadata]);

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
