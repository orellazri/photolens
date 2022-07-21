declare global {
  var API_URL: string;

  type Thumbnail = {
    id: number;
    createdAt: string;
    lastModified: string;
  };
}

export {};
