import { useState, useEffect } from "react";
import { Card, CardContent, CardMedia, Grid, Skeleton, Typography } from "@mui/material";
import axios from "axios";
import moment from "moment";

import "./style.css";

type GalleryProps = {
  limit?: Number;
  offset?: Number;
};

export default function Gallery({ limit, offset }: GalleryProps) {
  const [metadata, setMetadata] = useState<Array<Metadata>>([]);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);

  // Fetch metadata on page load
  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        const {
          data: { data },
        } = await axios.get(`/media/meta?limit=${limit}&offset=${offset}`);
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
    <Grid container spacing={1} className="grid">
      {metadata.map((_, i) =>
        thumbnails[i] ? (
          <Grid item key={i}>
            <a href={`${global.API_URL}/media/${thumbnails[i].id}`}>
              <Card>
                <CardMedia component="img" height="128" image={thumbnails[i].image} alt={thumbnails[i].id.toString()} />
                <CardContent>
                  <Typography sx={{ fontSize: 14 }} color="text.secondary" align="center" gutterBottom>
                    {thumbnails[i].createdAt}
                  </Typography>
                </CardContent>
              </Card>
            </a>
          </Grid>
        ) : (
          <Grid item key={i}>
            <Skeleton variant="rectangular" width={190} height={195} />
          </Grid>
        )
      )}
    </Grid>
  );
}
