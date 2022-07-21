import { useState, useEffect, useCallback } from "react";
import { FormControl, Grid, InputLabel, MenuItem, Select, SelectChangeEvent } from "@mui/material";
import { Box } from "@mui/system";
import axios from "axios";
import moment from "moment";

import "./style.css";
import ThumbnailCard from "../ThumbnailCard";

type SortByChoices = "created_at" | "last_modified";
type SortDirChoices = "desc" | "asc";
type Sort = {
  sortBy: SortByChoices;
  sortDir: SortDirChoices;
};

type ViewStyleChoices = "cards" | "tiles";

export default function Gallery() {
  const [isFetching, setIsFetching] = useState<boolean>(false);
  const [currentChunk, setCurrentChunk] = useState<number>(0);
  const [thumbnails, setThumbnails] = useState<Array<Thumbnail>>([]);
  const [sort, setSort] = useState<Sort>({ sortBy: "created_at", sortDir: "desc" });
  const [viewStyle, setViewStyle] = useState<ViewStyleChoices>("cards");

  const thumbnailsPerChunk = 50;

  // TODO: Add toasts to try catch blocks for errors

  // Load the next chunk
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
        .catch((e) => {});
    } catch (e) {
      console.error("Could not load chunk! " + e);
    }
  }, [currentChunk, sort]);

  // Handle the event where the user scroll to the bottom of the page
  const handleScrollToBottom = useCallback(() => {
    const bottom = Math.ceil(window.innerHeight + window.scrollY) >= document.documentElement.scrollHeight;
    if (!bottom) return;

    // Load the next chunk
    loadChunk();
  }, [loadChunk]);

  // Handle the event where the user changes the sort direction
  const handleChangeSortDir = (event: SelectChangeEvent) => {
    const eventData = (event.target.value as string).split("|");
    setSort({ sortBy: eventData[0] as SortByChoices, sortDir: eventData[1] as SortDirChoices });
    setThumbnails([]);
    setCurrentChunk(0);
  };

  // Handle the event where the user changes the view style
  const handleChangeViewStyle = (event: SelectChangeEvent) => {
    const eventData = event.target.value as ViewStyleChoices;
    setViewStyle(eventData);
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

  // Check whether we should load the next chunk (until the scroll bar appears
  // or we are out of photos to load)
  useEffect(() => {
    if (document.body.clientHeight <= window.innerHeight) {
      loadChunk();
    }
  }, [currentChunk, loadChunk]);

  // Render a thumbnail
  const renderThumbnail = (thumbnail: Thumbnail, i: number) => {
    switch (viewStyle) {
      case "cards":
        return (
          <Grid item key={i} xs={4} md={3} lg={2}>
            <ThumbnailCard thumbnail={thumbnail} />
          </Grid>
        );
      case "tiles":
        return (
          <Grid item key={i}>
            <a href={`${global.API_URL}/media/${thumbnail.id}`}>
              <img src={`${global.API_URL}/media/thumbnail/${thumbnail.id}`} alt={thumbnail.id.toString()} />
            </a>
          </Grid>
        );
    }
  };

  return (
    <Box>
      {/* Form */}
      <Box className="gallery-form">
        <FormControl disabled={isFetching} className="gallery-form-input">
          <InputLabel>Sort By</InputLabel>
          <Select value={`${sort.sortBy}|${sort.sortDir}`} label="Sort Direction" onChange={handleChangeSortDir}>
            <MenuItem value="created_at|desc">Recently Added</MenuItem>
            <MenuItem value="created_at|asc">Previously Added</MenuItem>
            <MenuItem value="last_modified|desc">Newest First</MenuItem>
            <MenuItem value="last_modified|asc">Oldest First</MenuItem>
          </Select>
        </FormControl>

        <FormControl className="gallery-form-input">
          <InputLabel>View Style</InputLabel>
          <Select value={viewStyle} label="View Style" onChange={handleChangeViewStyle}>
            <MenuItem value="cards">Cards</MenuItem>
            <MenuItem value="tiles">Tiles</MenuItem>
          </Select>
        </FormControl>
      </Box>

      {/* Grid */}
      <Grid container spacing={1} alignContent="center">
        {thumbnails.map((thumbnail, i) => renderThumbnail(thumbnail, i))}
      </Grid>
    </Box>
  );
}
