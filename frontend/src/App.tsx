import { useState, useEffect } from "react";
import axios from "axios";

function App() {
  useEffect(() => {
    const fetchMedia = async () => {
      let res = await axios.get("/media/2");
      console.log(res);
    };

    fetchMedia().catch(console.error);
  }, []);

  return (
    <div>
      <h1>Photolens</h1>
      <hr />
    </div>
  );
}

export default App;
