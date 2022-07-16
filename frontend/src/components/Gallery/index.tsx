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
  Tooltip,
  Typography,
} from "@mui/material";
import { AccessTime, Edit } from "@mui/icons-material";
import axios from "axios";
import moment from "moment";

import "./style.css";
import { Box } from "@mui/system";

type GalleryProps = {
  limit?: Number;
  offset?: Number;
};

type Sort = {
  sortBy: string;
  sortDir: string;
};

export default function Gallery({ limit = 0, offset = 0 }: GalleryProps) {
  const [isFetching, setIsFetching] = useState<boolean>(true);
  const [metadata, setMetadata] = useState<Array<Metadata>>([]);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);
  const [sort, setSort] = useState<Sort>({ sortBy: "created_at", sortDir: "desc" });

  // TODO: Fetch in chunks (configurable with prop) instead of all and then single requests

  // Fetch metadata on page load
  useEffect(() => {
    const fetchMetadata = async () => {
      try {
        setMetadata([]);
        setIsFetching(true);

        const {
          data: { data },
        } = await axios.get(`/media/meta?limit=${limit}&offset=${offset}&sortby=${sort.sortBy}&sortdir=${sort.sortDir}`);
        let metadataResults: Array<Metadata> = [];
        for (const result of data) {
          metadataResults.push({ id: result });
        }
        setMetadata(metadataResults);
      } catch (e) {
        console.error("Could not fetch metadata! " + e);
      }
    };

    fetchMetadata();
  }, [limit, offset, sort]);

  // Fetch thumbnails after fetching
  useEffect(() => {
    const fetchThumbnails = async () => {
      try {
        setThumbnails([]);

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
              lastModified: moment(data.last_modified).local().format("DD/MM/YYYY HH:mm:ss"),
            },
          ]);
        }
      } catch (e) {
        console.error("Could not fetch thumbnails! " + e);
      }
    };

    fetchThumbnails();
  }, [metadata]);

  useEffect(() => {
    if (metadata.length > 0 && thumbnails.length == metadata.length) {
      setIsFetching(false);
    }
  }, [thumbnails]);

  const handleChangeSortDir = (event: SelectChangeEvent) => {
    const eventData = (event.target.value as string).split("|");
    setSort({ sortBy: eventData[0], sortDir: eventData[1] });
  };

  return (
    <Box>
      {/* Form */}
      <Box className="form">
        <FormControl disabled={isFetching}>
          <InputLabel>Sort Direction</InputLabel>
          <Select value={`${sort.sortBy}|${sort.sortDir}`} label="Sort Direction" onChange={handleChangeSortDir}>
            <MenuItem value="created_at|desc">Recently Added</MenuItem>
            <MenuItem value="created_at|asc">Previously Added</MenuItem>
            <MenuItem value="last_modified|desc">Newest First</MenuItem>
            <MenuItem value="last_modified|asc">Oldest First</MenuItem>
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
                    <Tooltip title="Created At">
                      <Typography sx={{ fontSize: 14 }} color="text.secondary" className="label-with-icon" gutterBottom>
                        <AccessTime sx={{ fontSize: 14 }} />
                        &nbsp;
                        {thumbnails[i].createdAt}
                      </Typography>
                    </Tooltip>

                    <Tooltip title="Modified At">
                      <Typography sx={{ fontSize: 14 }} color="text.secondary" className="label-with-icon" gutterBottom>
                        <Edit sx={{ fontSize: 14 }} />
                        &nbsp;
                        {thumbnails[i].lastModified}
                      </Typography>
                    </Tooltip>
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
