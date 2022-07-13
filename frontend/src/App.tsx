import { useState, useEffect } from "react";
import axios from "axios";

function App() {
  const [thumbnails, setThumbnails] = useState<Array<string>>([]);

  useEffect(() => {
    const fetchThumbnails = async () => {
      const res = await axios.get("/media/thumbnail/all");
      for (let thumbnail of res.data.thumbnails) {
        setThumbnails((images) => [...images, "data:image/png;base64," + thumbnail]);
      }
    };

    fetchThumbnails();
  }, []);

  return (
    <div>
      <h1>Photolens</h1>
      <hr />
      {thumbnails.map((image) => (
        <img src={image} key={image} />
      ))}
    </div>
  );
}

export default App;
