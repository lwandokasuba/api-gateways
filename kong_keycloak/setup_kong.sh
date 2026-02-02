#!/bin/bash

KONG_ADMIN="http://localhost:8001"

echo "Waiting for Kong to be ready..."
until curl -s $KONG_ADMIN; do
  sleep 2
done

# --- Proxy Keycloak ---
echo "Configuring Keycloak Proxy..."
# Check if service exists (idempotent-ish check not strict, POST will fail if exists, which is fine)
# But strictly, we should check status. Assuming fail on conflict is okay or we can use PUT if we knew ID.
# Since this is a setup script, we assume clean or acceptable conflicts.

curl -i -X POST $KONG_ADMIN/services \
  --data name=keycloak-gateway \
  --data url=http://keycloak:8080

curl -i -X POST $KONG_ADMIN/services/keycloak-gateway/routes \
  --data name=keycloak-route \
  --data paths[]=/realms \
  --data paths[]=/resources \
  --data paths[]=/js \
  --data paths[]=/robots.txt

# --- Public Service ---
echo "Configuring Public Service..."
curl -i -X POST $KONG_ADMIN/services \
  --data name=service-public \
  --data url=http://service-public:8080

curl -i -X POST $KONG_ADMIN/services/service-public/routes \
  --data name=public-route \
  --data paths[]=/public

echo "Enabling Key Auth on Public Service..."
curl -i -X POST $KONG_ADMIN/services/service-public/plugins \
  --data name=key-auth

echo "Creating Consumer and Key..."
curl -i -X POST $KONG_ADMIN/consumers \
  --data username=external-user

curl -i -X POST $KONG_ADMIN/consumers/external-user/key-auth \
  --data key=apikey123

# --- Private Services ---
echo "Configuring Private Service 1..."
curl -i -X POST $KONG_ADMIN/services \
  --data name=service-private-1 \
  --data url=http://service-private-1:8080

curl -i -X POST $KONG_ADMIN/services/service-private-1/routes \
  --data name=private-1-route \
  --data paths[]=/private1

echo "Enabling OIDC on Private Service 1..."
# Point discovery to internal Keycloak.
# Keycloak (KC_HOSTNAME_URL=localhost:8000) will return issuer=localhost:8000.
# Kong validation passes as discovery url domain is ignored or irrelevant if signatures match.
# Wait, discovery URL is http://keycloak:8080. Body iss is http://localhost:8000.
# lua-resty-openidc might check this. If it fails, we set config.issuer explicitly.

curl -i -X POST $KONG_ADMIN/services/service-private-1/plugins \
  --data name=oidc \
  --data config.client_id=kong-client \
  --data config.client_secret=client-secret \
  --data config.discovery=http://localhost:8000/realms/kong-realm/.well-known/openid-configuration \
  --data config.logout_path=/logout \
  --data config.redirect_after_logout_uri=/

echo "Configuring Private Service 1 (Token/Bearer Only)..."
curl -i -X POST $KONG_ADMIN/services \
  --data name=service-private-1-token \
  --data url=http://service-private-1:8080

curl -i -X POST $KONG_ADMIN/services/service-private-1-token/routes \
  --data name=private-1-token-route \
  --data paths[]=/private-token

echo "Enabling OIDC (Bearer Only) on Private Service 1 Token..."
curl -i -X POST $KONG_ADMIN/services/service-private-1-token/plugins \
  --data name=oidc \
  --data config.client_id=kong-client \
  --data config.client_secret=client-secret \
  --data config.discovery=http://localhost:8000/realms/kong-realm/.well-known/openid-configuration \
  --data config.bearer_only=yes

echo "Configuring Private Service 2..."
curl -i -X POST $KONG_ADMIN/services \
  --data name=service-private-2 \
  --data url=http://service-private-2:8080

curl -i -X POST $KONG_ADMIN/services/service-private-2/routes \
  --data name=private-2-route \
  --data paths[]=/private2

curl -i -X POST $KONG_ADMIN/services/service-private-2/plugins \
  --data name=oidc \
  --data config.client_id=kong-client \
  --data config.client_secret=client-secret \
  --data config.discovery=http://localhost:8000/realms/kong-realm/.well-known/openid-configuration \
  --data config.logout_path=/logout \
  --data config.redirect_after_logout_uri=/

echo "Done!"
