
This repo sets up Kong with Postgres and two tiny Go services so you can test routing, plugins, and the Admin API without port conflicts.

## Ports (chosen to avoid common defaults)
- Postgres: `55432` -> `postgres:5432`
- Kong proxy: `49000` (http), `49443` (https)
- Kong Admin API: `49001` (http), `49444` (https)
- Kong Manager UI: `49002` (http), `49445` (https)
- Prometheus: `49090`
- Grafana: `49030`
- Go services (direct access):
  - svc-a: `49010`
  - svc-b: `49011`

## Prerequisites / knowledge
- Docker + Docker Compose basics
- Basic HTTP concepts (requests, routes, upstreams)
- `curl` or similar HTTP client
- (Optional) Kong concepts: services, routes, plugins, upstreams/targets

No local Go install required; the Go services run inside containers.

## Start everything
```sh
docker compose up -d postgres kong-migrations kong svc-a svc-b prometheus grafana
```

Check containers:
```sh
docker compose ps
```

## Verify the Go services directly (bypass Kong)
```sh
curl http://localhost:49010/
curl http://localhost:49011/
```

## Configure Kong (Admin API)
Create two services and routes:
```sh
# svc-a
curl -sS -X POST http://localhost:49001/services \
  --data name=svc-a \
  --data url=http://svc-a:8080

curl -sS -X POST http://localhost:49001/routes \
  --data name=svc-a-route \
  --data service.name=svc-a \
  --data 'paths[]=/a'

# svc-b
curl -sS -X POST http://localhost:49001/services \
  --data name=svc-b \
  --data url=http://svc-b:8080

curl -sS -X POST http://localhost:49001/routes \
  --data name=svc-b-route \
  --data service.name=svc-b \
  --data 'paths[]=/b'
```

## Test through Kong (proxy)
```sh
curl http://localhost:49000/a
curl http://localhost:49000/b
```

## Kong Manager UI
Open: `http://localhost:49002`

Note: Kong Manager is an Enterprise feature. With the OSS image, the UI shell may load but won’t allow management.
Use the Admin API (curl) for OSS.

## Suggested experiments
- Add a rate-limiting plugin to one route
- Enable key-auth and test with/without keys
- Create an upstream + targets and load-balance across multiple instances

## Metrics and monitoring (Prometheus + Grafana)
Enable the Prometheus plugin globally:
```sh
curl -sS -X POST http://localhost:49001/plugins \
  --data name=prometheus
```

Verify Kong metrics:
```sh
curl -sS http://localhost:49001/metrics | head -n 20
```

### Grafana setup
Open Grafana: `http://localhost:49030` (login `admin` / `admin`, then change password)

Add Prometheus datasource:
- Connections → Data sources → Add data source → Prometheus
- URL: `http://prometheus:9090`
- Save & test

Import a Kong dashboard:
- Dashboards → New → Import
- Search for an official Kong dashboard in Grafana’s catalog and import it
- Select the Prometheus datasource you added

## Stop and clean up
```sh
docker compose down -v
```
