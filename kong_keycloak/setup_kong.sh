#!/bin/bash

KONG_ADMIN="http://localhost:8001"

echo "Waiting for Kong to be ready..."
until curl -s $KONG_ADMIN; do
  sleep 2
done

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

echo "Configuring Private Service 1..."
curl -i -X POST $KONG_ADMIN/services \
  --data name=service-private-1 \
  --data url=http://service-private-1:8080

curl -i -X POST $KONG_ADMIN/services/service-private-1/routes \
  --data name=private-1-route \
  --data paths[]=/private1

echo "Enabling OIDC on Private Service 1..."
# Note: discovery url might need to use internal docker network name if kong -> keycloak
# but for browser redirection, the issuer in token must match what browser sees.
# This is tricky with Docker networking + localhost.
# Keycloak frontend URL needs to be set if accessed from host.
# For simplicity, we assume everything is localhost for now, but Kong container needs to resolve 'keycloak'.
# We often need to set KC_HOSTNAME or KC_HOSTNAME_URL.

curl -i -X POST $KONG_ADMIN/services/service-private-1/plugins \
  --data name=oidc \
  --data config.client_id=kong-client \
  --data config.client_secret=client-secret \
  --data config.discovery=http://keycloak:8080/realms/kong-realm/.well-known/openid-configuration \
  --data config.redirect_uri_path_pattern="/private1"
# For OIDC to work comfortably in docker, ensure 'keycloak' hostname is resolvable by Kong.
# And browser can also see 'keycloak' (via /etc/hosts) OR use localhost throughout.
# If using localhost for discovery, Kong must be able to hit localhost:8080 (which is itself? No. It's host).
# Best is to use internal docker name for discovery: http://keycloak:8080/...

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
  --data config.discovery=http://keycloak:8080/realms/kong-realm/.well-known/openid-configuration \
  --data config.redirect_uri_path_pattern="/private2"

echo "Done!"
