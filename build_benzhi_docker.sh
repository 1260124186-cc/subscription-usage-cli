#!/usr/bin/env sh
set -eu

platform="${1:?usage: ./build_benzhi_docker.sh linux/amd64|linux/arm64}"
tag="subscription-usage-cli:benzhi-$(printf '%s' "$platform" | tr '/' '-')"

docker buildx build --platform "$platform" --load --file benzhi.Dockerfile --tag "$tag" .
docker run --rm "$tag" go build ./...
docker run --rm "$tag"
