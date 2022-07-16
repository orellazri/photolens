import { useState, useEffect } from "react";
import {
  Card,
  CardContent,
  CardMedia,
  FormControl,
  Grid,
  InputLabel,
  MenuItem,
  Select,
  SelectChangeEvent,
  Skeleton,
  Typography,
} from "@mui/material";
import axios from "axios";
import moment from "moment";

import "./style.css";
import { Box } from "@mui/system";

type GalleryProps = {
  limit?: Number;
  offset?: Number;
};

export default function Gallery({ limit = 0, offset = 0 }: GalleryProps) {
  const [metadata, setMetadata] = useState<Array<Metadata>>([]);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);
  const [sortDir, setSortDir] = useState<string>("desc");

  // TODO: Fetch in chunks (configurable with prop)

  // Fetch metadata on page load
  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        const {
          data: { data },
        } = await axios.get(`/media/meta?limit=${limit}&offset=${offset}&sortdir=${sortDir}`);
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

    setMetadata([]);
    fetchMetadata();
  }, [limit, offset, sortDir]);

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

    setThumbnails([]);
    fetchThumbnails();
  }, [metadata]);

  const handleChangeSortDir = (event: SelectChangeEvent) => {
    setSortDir(event.target.value as string);
  };

  return (
    <Box>
      <Box className="form">
        <FormControl>
          <InputLabel>Sort Direction</InputLabel>
          <Select value={sortDir} label="Sort Direction" onChange={handleChangeSortDir}>
            <MenuItem value="desc">Descending</MenuItem>
            <MenuItem value="asc">Ascending</MenuItem>
          </Select>
        </FormControl>
      </Box>

      {/* Grid */}
      <Grid container spacing={1} className="grid">
        {metadata.map((_, i) =>
          thumbnails[i] ? (
            // Thumbnails
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
            // Placeholders
            <Grid item key={i}>
              <Skeleton variant="rectangular" width={190} height={195} />
            </Grid>
          )
        )}
      </Grid>
    </Box>
  );
}
