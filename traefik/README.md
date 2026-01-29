# Traefik + Go Microservices Demo

This repo runs Traefik as a reverse proxy with three small Go services (`alpha`, `beta`, `gamma`) and a gateway tester that generates both successful and error traffic.

## Prereqs
- Docker + Docker Compose
- Go (only needed to run the gateway test client)

## Run everything
```bash
docker compose up --build
```

Traefik dashboard (dev-only):
- http://localhost:8080/

## Test via Traefik (recommended)
Run the Go test client, which sends requests through Traefik using Host headers and deliberately generates errors.

```bash
go run gateway_test.go
```

What it hits per service:
- `/` (200)
- `/health` (200)
- `/delay` (200, slower)
- `/error` (500)
- `/missing` (404)

You should see the errors and latency reflected in the Traefik dashboard.

## Manual curl (optional)
```bash
curl -H 'Host: alpha.localhost' http://localhost/
curl -H 'Host: beta.localhost' http://localhost/error
curl -H 'Host: gamma.localhost' http://localhost/missing
```

## Services and routes
- `alpha.localhost`
- `beta.localhost`
- `gamma.localhost`
- `whoami.localhost`

## Notes
- Traefik API/dashboard is exposed insecurely on port 8080 for local dev only.
- If you see a network error, ensure Docker is running and the compose stack is up.
