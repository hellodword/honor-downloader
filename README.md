# honor-downloader

## build

```sh
env GOWORK=off GOTOOLCHAIN=local CGO_ENABLED=0 go build -x -v -trimpath -ldflags "-s -w" -buildvcs=false -o ./dist/honor-downloader .

# or using docker
docker build -t honor-downloader .
```

## usage

```sh
# if using docker
alias honor-downloader='docker run --rm --name honor-downloader -v "$(pwd)":/tmp --user "$(id -u):$(id -g)" honor-downloader'

honor-downloader -anna https://annas-archive.li/md5/08b0f97b98c977da93cd5e5623686af5 -name "The Embodied Soul: Aristotelian Psychology and Physiology in Medieval Europe between 1200 and 1420.epub"
```

## unsupported

- https://annas-archive.se/md5/c8acc9e8b9ccfb34f9630cab84c8060b
