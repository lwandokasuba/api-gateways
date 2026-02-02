# Kong + Keycloak + Go Services + Custom Dashboard

This project provides a complete setup for an API Gateway solution using Kong and Keycloak, with example Go microservices and a custom React Dashboard.

## 🚀 Features

- **Kong Gateway** (with `kong-oidc` plugin installed).
- **Keycloak** Identity Provider (pre-configured via import).
- **3 Go Microservices**:
  - `service-public`: Accessible via API Key.
  - `service-private-1`: Secured via OIDC (Keycloak).
  - `service-private-2`: Secured via OIDC (Keycloak).
- **Custom Dashboard**: React + Vite application to monitor services.
- **Docker Compose** orchestration.

## 📂 Directory Structure

- `docker-compose.yml`: Orchestrates Kong, Keycloak, Database, Services, and Dashboard (dev).
- `kong/`: Custom Kong Dockerfile.
- `keycloak-import/`: Realm export for bootstrapping Keycloak.
- `services/`: Source code for Go microservices.
- `dashboard/`: React Dashboard source code.
- `setup_kong.sh`: Script to configure Kong routes and plugins.

## 🛠 Setup Instructions

### 1. Start Infrastructure
Run the following command to build and start all services:

```bash
docker-compose up --build -d
```
Wait for all containers to be healthy. Keycloak and Kong migrations might take a minute.

### 2. Configure Kong
Once Kong is running (Admin API at localhost:8001), run the setup script:

```bash
./setup_kong.sh
```
This script will:
- register the Go services in Kong.
- create Routes (`/public`, `/private1`, `/private2`).
- enable `key-auth` for public route.
- enable `oidc` for private routes (using `kong-oidc` plugin).
- create a consumer `external-user` with api key `apikey123`.

### 3. Run Dashboard
The dashboard is a separate frontend application.

```bash
cd dashboard
npm install
npm run dev
```
Open [http://localhost:5173](http://localhost:5173) to view the dashboard.

## 🧪 Verification

### Public API (Key Auth)
**Without Key:**
```bash
curl -i http://localhost:8000/public
# HTTP/1.1 401 Unauthorized
```

**With Key:**
```bash
curl -i -H "apikey: apikey123" http://localhost:8000/public
# HTTP/1.1 200 OK
```

### Private API (OIDC)
**Without Token:**
```bash
curl -i http://localhost:8000/private1
# HTTP/1.1 401 Unauthorized (or Redirect to Keycloak)
```

**Get Token & Access:**
You can use the browser to visit `http://localhost:8000/private1`. You will be redirected to Keycloak.
* Credentials: `user1` / `password`.
* After login, you will see the JSON response from the service.

### Private API (Bearer Token Route)
This route `http://localhost:8000/private-token` is configured with `bearer_only=yes`, meaning it will not redirect to login but expects a valid Bearer token.

**1. Get Token:**
```bash
export TOKEN=$(curl -s -X POST http://localhost:8000/realms/kong-realm/protocol/openid-connect/token \
  -d client_id=kong-client \
  -d client_secret=client-secret \
  -d username=user1 \
  -d password=password \
  -d grant_type=password | jq -r .access_token)
echo $TOKEN
```

**2. Access Service:**
```bash
curl -i -H "Authorization: Bearer $TOKEN" http://localhost:8000/private-token
# HTTP/1.1 200 OK
```

## 🔐 Keycloak Details
- **Console**: [http://localhost:8080](http://localhost:8080)
- **Admin**: `admin` / `admin`
- **Realm**: `kong-realm` (Imported automatically)
- **Client**: `kong-client`
- **User**: `user1` / `password`

## 📊 Dashboard
The dashboard uses the Kong Admin API to display status:
- View active Services.
- View Routes.
- Monitor basic stats (mocked/API).

---
**Note:** Ensure `kong-oidc` plugin is correctly loaded. The custom Dockerfile in `kong/` handles this.
