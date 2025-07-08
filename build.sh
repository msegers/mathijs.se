#!/bin/sh
# build.sh
set -euo pipefail

set -e
INCREMENT_MODE="$1"
CHART_FILE="./helm/Chart.yaml"

if [[ $INCREMENT_MODE = '' ]] ; then
  INCREMENT_MODE='patch'
fi

CURRENT_VERSION=$(grep -Eo "([0-9]+\.[0-9]+\.[0-9]+)" "$CHART_FILE" --max-count 1)

SEGMENTS=(${CURRENT_VERSION//./ });
MAJOR=${SEGMENTS[0]}
MINOR=${SEGMENTS[1]}
PATCH=${SEGMENTS[2]}

if [ "$INCREMENT_MODE" = 'major' ]; then
  MAJOR=$((MAJOR+1))
  MINOR="0"
  PATCH="0"
elif [ "$INCREMENT_MODE" = 'minor' ]; then
  MINOR=$((MINOR+1))
  PATCH="0"
elif [ "$INCREMENT_MODE" = 'patch' ]; then
  PATCH=$((PATCH+1))
else
  echo 'USAGE: build.sh ['patch' | 'minor' | 'major']'
  echo 'Please provide a proper target semver increment (leaving it empty uses patch)'
  exit 0
fi

NEW_VERSION="$MAJOR.$MINOR.$PATCH"

echo $NEW_VERSION

# compiles the app
hugo

# make docker image and push
sudo docker buildx build --platform linux/arm64,linux/amd64 -t mathijs-se .
sudo docker tag mathijs-se msegers/mathijs-se:"$NEW_VERSION"
sudo docker push msegers/mathijs-se:"$NEW_VERSION"

# update helm file
sed -i -e "s/$CURRENT_VERSION/$NEW_VERSION/g" "$CHART_FILE"
exit 0;

