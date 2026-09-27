import React, {useState} from "react";
import { shortenUrl, getStats } from "./api";
import './App.css';

function App() {
  const [url,setUrl] = useState("");
  const [shortUrl, setShortUrl] = useState("");
  const [clicks, setClicks] = useState(null)
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleShorten = async () => {
    setError("");
    setShortUrl("");
    setClicks(null);
    setLoading(true);
    try {
      const res = await shortenUrl(url);
      setShortUrl(res.short_url);

      const shortCode = res.short_url.split("/").pop();
      const stats = await getStats(shortCode);
      setClicks(stats.clicks)
    } catch (err) {
      setError("URL kısaltılamadı. Adresi ve API bağlantısını kontrol edin.");
    } finally {
      setLoading(false);
    }
  }

  
  return (
    <div style={{ padding: 40 }}>
      <h2>URL Kısaltıcı</h2>
      <input
        type="text"
        aria-label="Kısaltılacak URL"
        value={url}
        onChange={(e) => setUrl(e.target.value)}
        placeholder="URL giriniz"
        style={{ width: "300px", marginRight: 10 }}
      />
      <button onClick={handleShorten} disabled={loading || !url.trim()}>{loading ? "Kısaltılıyor..." : "Shorten"}</button>
      {error && <p role="alert">{error}</p>}

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






