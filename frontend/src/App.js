import React, {useState} from "react";
import { shortenUrl, getStats } from "./api";
import './App.css';

function App() {
  const [url,setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [clicks, setClicks] = useState(null)

  const handleShorten = async () => {
    try {
      const res = await shortenUrl(url);
      setShortUrl(res.short_url);

      const shortCode = res.short_url.split("/").pop();
      const stats = await getStats(shortCode);
      setClicks(stats.clicks)
    } catch (err) {
      
    }
  }

  
  return (
    <div style={{ padding: 40 }}>
      <h2>URL Kısaltıcı</h2>
      <input
        type="text"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="URL giriniz"
        style={{ width: "300px", marginRight: 10 }}
      />
      <button onClick={handleShorten}>Shorten</button>

      {shortUrl && (
        <div style={{ marginTop: 20 }}>
          <p>
            <strong>Kısaltılmış URL:</strong>{" "}
            <a href={shortUrl} target="_blank" rel="noreferrer">
              {shortUrl}
            </a>
          </p>
          <p><strong>Tıklanma Sayısı:</strong> {clicks}</p>
        </div>
      )}
    </div>
  );
}


export default App;






