FROM debian:bookworm

RUN apt-get update && \
    DEBIAN_FRONTEND="noninteractive" \
    apt-get install -y \
    curl

RUN LATEST_VERSION="$(curl -fsS -w "%{redirect_url}" -o /dev/null "https://github.com/ikatson/rqbit/releases/latest" | grep -oP '(?<=/releases/tag/)[^/]+$')" && \
    curl -fSL --output /usr/local/bin/rqbit "https://github.com/ikatson/rqbit/releases/download/$LATEST_VERSION/rqbit-linux-static-x86_64" && \
    chmod +x /usr/local/bin/rqbit

COPY honor-downloader /usr/local/bin/honor-downloader
COPY entrypoint.sh /entrypoint.sh

RUN chmod +x /usr/local/bin/honor-downloader && chmod +x /entrypoint.sh

WORKDIR /data

ENTRYPOINT ["/entrypoint.sh"]
