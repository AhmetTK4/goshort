# GoShort - Simple URL Shortener with Go and Redis

GoShort is a minimal and fast URL shortener service written in Go using the Gin framework and Redis for storage.

## Features

- Shorten long URLs
- Redirect short URLs to original ones
- Track how many times a short URL was clicked
- Docker + Docker Compose ready

## Tech Stack

- Go 1.24.5+
- Gin Web Framework
- Redis (as key-value store)
- Docker + Docker Compose

## Running Locally

### 1. Clone the repo

```bash
git clone https://github.com/AhmetTK4/goshort.git
cd goshort
```

### 2. Build and run with Docker

```bash
docker compose up --build
```

- API runs at: http://localhost:8080
- Redis is reachable by the API at `redis:6379` inside Compose; its port is not published to the host.

### 3. Run the React interface

In another terminal, with Node.js 22 and npm installed:

```sh
cd frontend
npm ci
npm start
```

Open `http://localhost:3000`, paste an HTTP(S) URL, and click **Shorten**. Open the resulting link to trigger a redirect. The interface fetches the click count when shortening; submit the same original URL again to refresh its count. Compose starts the API and Redis, not the frontend.

### Configuration

- `REDIS_ADDR`: API's Redis address; defaults to `localhost:6379` outside Docker. Compose sets `redis:6379`.
- `BASE_URL`: externally reachable API URL used to create short links; defaults to `http://localhost:8080`.
- `REACT_APP_API_BASE_URL`: frontend API URL, read when starting/building the frontend; defaults to `http://localhost:8080`.

The API currently permits browser requests from `http://localhost:3000`. Configure an explicit allowed origin before hosting the frontend elsewhere. Only absolute HTTP(S) URLs without embedded credentials are accepted.

## Verification

```sh
go test ./...
go vet ./...
cd frontend
npm ci
npm test -- --watchAll=false
npm run build
```

API tests cover URL creation, redirects, click counts, invalid URLs, and unknown codes with an in-memory Redis-compatible test server. Frontend tests cover empty input, successful shortening, and API failures. See [CONTRIBUTING.md](CONTRIBUTING.md).

## Scope and limitations

This is a local learning project. Links expire after one hour; click counters currently do not share that expiry. Redis data is not configured for durable storage. Rate limiting, abuse prevention, collision handling, and hardened deployment are not implemented. Do not expose the demo as a public URL-shortening service without addressing those concerns. No performance benchmark is claimed.

## API Endpoints

### POST /api/shorten

Shortens a long URL.

Request:
```json
{
  "url": "https://example.com/very/long/url"
}
```

Response:
```json
{
  "short_url": "http://localhost:8080/g/AbC123"
}
```

### GET /g/:shortCode

Redirects to original URL.

Example:
```
http://localhost:8080/g/AbC123
```

### GET /api/stats/:shortCode

Returns how many times the short URL was clicked.

Response:
```json
{
  "short_code": "AbC123",
  "clicks": "5"
}
```

## Testing with curl

```bash
curl -X POST http://localhost:8080/api/shorten \
-H "Content-Type: application/json" \
-d '{"url": "https://openai.com"}'
```

## Project Structure

```
.
├── main.go
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
├── /service
└── /storage
```

## Author

Made by Ahmet Temel Kundupoğlu - https://github.com/AhmetTK4
