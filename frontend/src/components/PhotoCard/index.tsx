import { Card, CardActionArea, CardContent, CardMedia, Tooltip, Typography } from "@mui/material";
import { AccessTime, Edit } from "@mui/icons-material";

import "./style.css";

type PhotoCardProps = {
  thumbnail: Thumbnail;
};

export default function PhotoCard({ thumbnail }: PhotoCardProps) {
  return (
    <Card>
      <CardActionArea>
        <a href={`${global.API_URL}/media/${thumbnail.id}`}>
          <CardMedia component="img" height="128" image={thumbnail.image} alt={thumbnail.id.toString()} />
          <CardContent>
            <Tooltip title="Created At">
              <Typography sx={{ fontSize: 14 }} color="text.secondary" className="label-with-icon" gutterBottom>
                <AccessTime sx={{ fontSize: 14 }} />
                &nbsp;
                {thumbnail.createdAt}
              </Typography>
            </Tooltip>

            <Tooltip title="Modified At">
              <Typography sx={{ fontSize: 14 }} color="text.secondary" className="label-with-icon" gutterBottom>
                <Edit sx={{ fontSize: 14 }} />
                &nbsp;
                {thumbnail.lastModified}
              </Typography>
            </Tooltip>
          </CardContent>
        </a>
      </CardActionArea>
    </Card>
  );
}
