import { useState, useEffect } from "react";
import { Buffer } from "buffer";
import axios from "axios";

function App() {
  const [image, setImage] = useState<string>("");

  useEffect(() => {
    const fetchMedia = async () => {
      let res = await axios.get("/media/thumbnail/2");
      setImage("data:image/png;base64," + res.data);
    };

    fetchMedia();
  }, []);

  return (
    <div>
      <h1>Photolens</h1>
      <hr />
      <img src={image} />
    </div>
  );
}

export default App;
