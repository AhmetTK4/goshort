# GoShort - Simple URL Shortener with Go and Redis

GoShort is a minimal and fast URL shortener service written in Go using the Gin framework and Redis for storage.

## Features

- Shorten long URLs
- Redirect short URLs to original ones
- Track how many times a short URL was clicked
- Docker + Docker Compose ready

## Tech Stack

- Go 1.24+
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
- Redis runs at: localhost:6379 (internal in Docker)

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
