FROM golang:bookworm AS builder

ENV GOWORK=off
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
COPY go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN commit="$(git rev-parse HEAD)" && go build -x -v -trimpath -ldflags "-s -w -X main.commit=${commit}" -buildvcs=false -o /usr/local/bin/honor-downloader .

FROM gcr.io/distroless/base-debian12

COPY --from=builder /usr/local/bin/honor-downloader /usr/local/bin/honor-downloader

WORKDIR /tmp

ENTRYPOINT ["/usr/local/bin/honor-downloader"]
