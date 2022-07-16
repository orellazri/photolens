declare global {
  var API_URL: string;

  type Thumbnail = {
    id: number;
    image: string;
    createdAt: string;
  };

  type Metadata = {
    id: number;
    createdAt: string;
  };
}

export {};
