import { useState, useEffect, useCallback } from "react";
import { FormControl, Grid, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";
import moment from "moment";

import "./style.css";
import PhotoCard from "../PhotoCard";

type Sort = {
  sortBy: string;
  sortDir: string;
};

export default function Gallery() {
  const [isFetching, setIsFetching] = useState<boolean>(false);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);
  const [sort, setSort] = useState<Sort>({ sortBy: "created_at", sortDir: "desc" });
  const [currentChunk, setCurrentChunk] = useState<number>(0);
  const [err, setErr] = useState<any>(null);

  const thumbnailsPerChunk = 50;

  // TODO: Add toasts to try catch blocks for errors

  const loadChunk = useCallback(() => {
    try {
      setIsFetching(true);

      axios
        .get(
          `/media/meta?limit=${thumbnailsPerChunk}&offset=${thumbnailsPerChunk * currentChunk}&sortby=${sort.sortBy}&sortdir=${
            sort.sortDir
          }`
        )
        .then((res) => {
          if (!res.data.data) {
            setIsFetching(false);
            return;
          }
          const {
            data: { data },
          } = res;
          let thumbnailsResults: Array<Thumbnail> = [];
          for (const result of data) {
            thumbnailsResults.push({
              id: result.id,
              createdAt: moment(result.created_at).local().format("DD/MM/YYYY"),
              lastModified: moment(result.last_modified).local().format("DD/MM/YYYY"),
            });
          }
          setThumbnails((thumbnails) => [...thumbnails, ...thumbnailsResults]);
          setCurrentChunk((currentChunk) => currentChunk + 1);
          setIsFetching(false);
        })
        .catch((e) => {
          setErr("---> 1" + e);
        });
    } catch (e) {
      console.error("Could not load chunk! " + e);
      setErr("---> 2" + e);
    }
  }, [currentChunk, sort]);

  const handleScrollToBottom = useCallback(() => {
    const bottom = Math.ceil(window.innerHeight + window.scrollY) >= document.documentElement.scrollHeight;
    if (!bottom) return;
    loadChunk();
  }, [loadChunk]);

  const handleChangeSortDir = (event: SelectChangeEvent) => {
    const eventData = (event.target.value as string).split("|");
    setSort({ sortBy: eventData[0], sortDir: eventData[1] });
    setThumbnails([]);
    setCurrentChunk(0);
  };

  useEffect(() => {
    // Scroll listener
    window.addEventListener("scroll", handleScrollToBottom, {
      passive: true,
    });

    return () => {
      window.removeEventListener("scroll", handleScrollToBottom);
    };
  }, [currentChunk, handleScrollToBottom]);

  useEffect(() => {
    if (document.body.clientHeight <= window.innerHeight) {
      loadChunk();
    }
  }, [currentChunk, loadChunk]);

  return (
    <Box>
      {err}
      {/* Form */}
      <Box className="form">
        <FormControl disabled={isFetching}>
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
