declare global {
  var API_URL: string;

  type Thumbnail = {
    id: number;
    image: string;
    createdAt: string;
    lastModified: string;
  };

  type Metadata = {
    id: number;
  };
}

export {};
