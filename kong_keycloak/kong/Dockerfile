FROM kong:latest

USER root
# Install dependencies for luarocks
RUN apt-get update && apt-get install -y git unzip curl wget build-essential

RUN luarocks install lua-resty-openidc

# Install kong-oidc from revomatico fork for Kong 3.x support
RUN git clone https://github.com/revomatico/kong-oidc.git \
    && cd kong-oidc \
    && luarocks make

USER kong
