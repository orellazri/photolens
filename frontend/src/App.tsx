import { useState, useEffect } from "react";
import axios from "axios";

function App() {
  const [images, setImages] = useState<Array<string>>([]);

  useEffect(() => {
    const fetchMedia = async () => {
      const res = await axios.get("/media/thumbnail/all");
      console.log(res);

      for (let thumbnail of res.data.thumbnails) {
        setImages((images) => [...images, "data:image/png;base64," + thumbnail]);
      }
    };

    fetchMedia();
  }, []);

  return (
    <div>
      <h1>Photolens</h1>
      <hr />
      {images.map((image) => (
        <img src={image} />
      ))}
    </div>
  );
}

export default App;
