# vaqas.dev

My personal portfolio site, built with Go and deployed on a self-managed VPS. It displays NASA's Astronomy Picture of the Day alongside my personal links and projects.

Live at [vaqas.dev](https://vaqas.dev).

## What it does

- Serves dynamically rendered HTML templates using Go's `html/template`
- Fetches the daily image from NASA's APOD REST API
- Caches the API response until midnight EST so the API is only called once per day
- Serves static files for CSS and other assets
- Runs in a Docker container behind a Caddy reverse proxy with automatic TLS

## How the caching works

Every request checks an in-memory cache keyed by the current UTC date. If the cache has a valid entry, the stored image URL is used directly. If not, the server calls the NASA API, stores the result, and sets the cache expiry to exactly midnight EST. This aligns with when NASA publishes a new image each day.

Early on, the server was calling the NASA API on every request. This caused repeated latency and occasional failures when the API returned unexpected responses. The cache fixes both problems.

## What I learned

- How to structure a Go web server using `net/http` with clean handler functions
- How in-memory caching works and how to set expiry times tied to real-world events like a daily publish schedule
- How to handle external API failures gracefully without crashing the server
- How to containerize a Go app with Docker and wire it up to Caddy for TLS and reverse proxying
- How `log.Fatal` inside a handler kills the entire server process, and why error handling matters in long-running programs

## Stack

- Go (`net/http`, `html/template`, `encoding/json`)
- go-cache
- Docker, Docker Compose
- Caddy
- NASA APOD REST API

## Running locally

Requires a NASA API key. Get one free at [api.nasa.gov](https://api.nasa.gov).

```bash
git clone https://github.com/vaqasq/vaqas.dev
cd vaqas.dev
echo "NASA_API_KEY=your_key_here" > .env
go run main.go
```

Site available at `http://localhost:8080`.