#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 2 ]]; then
  echo "Usage: $0 /absolute/path/to/image content/blog/YYYY/MM/slug" >&2
  exit 1
fi

if ! command -v ffmpeg >/dev/null 2>&1; then
  echo "ffmpeg is required for image resizing and thumbnail generation." >&2
  exit 1
fi

SOURCE="$1"
POST_DIR="$2"

if [[ ! -f "$SOURCE" ]]; then
  echo "Image not found: $SOURCE" >&2
  exit 1
fi

EXTENSION="${SOURCE##*.}"
EXTENSION="${EXTENSION,,}"
case "$EXTENSION" in
  jpg|jpeg)
    EXTENSION="jpg"
    QUALITY=(-q:v 3)
    ;;
  png)
    QUALITY=(-compression_level 9)
    ;;
  webp)
    QUALITY=(-quality 80)
    ;;
  *)
    echo "Unsupported image type: .$EXTENSION (use jpg, png, or webp)" >&2
    exit 1
    ;;
esac

mkdir -p "$POST_DIR"

IMAGE="$POST_DIR/image.$EXTENSION"
THUMB="$POST_DIR/thumb.$EXTENSION"

# Keep the full image below 1600px on the longest side, without upscaling.
ffmpeg -hide_banner -loglevel error -y -i "$SOURCE" \
  -vf "scale='min(1600,iw)':'min(1600,ih)':force_original_aspect_ratio=decrease" \
  -frames:v 1 "${QUALITY[@]}" "$IMAGE"

# A 16:9 thumbnail for the blog overview cards.
ffmpeg -hide_banner -loglevel error -y -i "$SOURCE" \
  -vf "scale=800:450:force_original_aspect_ratio=increase,crop=800:450" \
  -frames:v 1 "${QUALITY[@]}" "$THUMB"

echo "image: $IMAGE ($(du -h "$IMAGE" | cut -f1))"
echo "thumb: $THUMB ($(du -h "$THUMB" | cut -f1))"
