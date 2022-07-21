import { useState, useEffect } from "react";
import { FormControl, Grid, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";
import moment from "moment";

import "./style.css";
import PhotoCard from "../PhotoCard";

type GalleryProps = {
  limit?: Number;
  offset?: Number;
};

type Sort = {
  sortBy: string;
  sortDir: string;
};

export default function Gallery({ limit = 0, offset = 0 }: GalleryProps) {
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);
  const [sort, setSort] = useState<Sort>({ sortBy: "created_at", sortDir: "desc" });

  // TODO: Fetch in chunks (configurable with prop) instead of all and then single requests
  // TODO: Add toasts to try catch blocks for errors

  useEffect(() => {
    const fetchThumbnails = async () => {
      try {
        setThumbnails([]);

        const {
          data: { data },
        } = await axios.get(`/media/meta?limit=${limit}&offset=${offset}&sortby=${sort.sortBy}&sortdir=${sort.sortDir}`);
        let thumbnailsResults: Array<Thumbnail> = [];
        for (const result of data) {
          thumbnailsResults.push({
            id: result.id,
            createdAt: moment(result.created_at).local().format("DD/MM/YYYY"),
            lastModified: moment(result.last_modified).local().format("DD/MM/YYYY"),
          });
        }
        setThumbnails(thumbnailsResults);
      } catch (e) {
        console.error("Could not fetch metadata! " + e);
      }
    };

    fetchThumbnails();
  }, [limit, offset, sort]);

  const handleChangeSortDir = (event: SelectChangeEvent) => {
    const eventData = (event.target.value as string).split("|");
    setSort({ sortBy: eventData[0], sortDir: eventData[1] });
  };

  return (
    <Box>
      {/* Form */}
      <Box className="form">
        <FormControl>
          <InputLabel>Sort By</InputLabel>
          <Select value={`${sort.sortBy}|${sort.sortDir}`} label="Sort Direction" onChange={handleChangeSortDir}>
            <MenuItem value="created_at|desc">Recently Added</MenuItem>
            <MenuItem value="created_at|asc">Previously Added</MenuItem>
            <MenuItem value="last_modified|desc">Newest First</MenuItem>
            <MenuItem value="last_modified|asc">Oldest First</MenuItem>
          </Select>
        </FormControl>
      </Box>

      {/* Grid */}
      <Grid container spacing={1} alignContent="center">
        {thumbnails.map((thumbnail, i) => (
          <Grid xs={4} md={3} lg={2} item key={i}>
            <PhotoCard thumbnail={thumbnail} />
          </Grid>
        ))}
      </Grid>
    </Box>
  );
}
