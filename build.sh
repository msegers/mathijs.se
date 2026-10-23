#!/usr/bin/env bash
set -euo pipefail

MODE="${1:-patch}"
CHART_FILE="./helm/Chart.yaml"
IMAGE="msegers/mathijs-se"

if [[ ! -f "$CHART_FILE" ]]; then
  echo "Missing $CHART_FILE" >&2
  exit 1
fi

CURRENT_VERSION="$(grep -Eo '([0-9]+\.[0-9]+\.[0-9]+)' "$CHART_FILE" | head -n1)"
IFS=. read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

case "$MODE" in
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  patch)
    PATCH=$((PATCH + 1))
    ;;
  *)
    echo "Usage: ./build.sh [patch|minor|major]" >&2
    exit 1
    ;;
esac

NEW_VERSION="$MAJOR.$MINOR.$PATCH"
echo "Building version $NEW_VERSION"

# Build the static site locally first, so content/template errors fail fast.
go run .

# Multi-platform builds need a container-based builder; the default "docker"
# driver cannot produce a multi-arch manifest.
BUILDER="${BUILDX_BUILDER:-mathijs-se-builder}"
if ! docker buildx inspect "$BUILDER" >/dev/null 2>&1; then
  echo "Creating buildx builder: $BUILDER"
  docker buildx create --name "$BUILDER" --driver docker-container --bootstrap
fi

if [[ ! -f "${DOCKER_CONFIG:-$HOME/.docker}/config.json" ]]; then
  echo "Note: no Docker config found. Run 'docker login' if the push is rejected." >&2
fi

# buildx can build both architectures and push the result directly.
docker buildx build \
  --builder "$BUILDER" \
  --platform linux/arm64,linux/amd64 \
  -t "$IMAGE:$NEW_VERSION" \
  --push .

# Keep the Helm chart version in sync with the image tag.
ESCAPED_VERSION="${CURRENT_VERSION//./\\.}"
sed -i \
  -e "s/^version: $ESCAPED_VERSION$/version: $NEW_VERSION/" \
  -e "s/^appVersion: $ESCAPED_VERSION$/appVersion: $NEW_VERSION/" \
  "$CHART_FILE"

echo "Updated $CHART_FILE to $NEW_VERSION"
