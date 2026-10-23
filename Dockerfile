# syntax=docker/dockerfile:1

# The site generator only needs to run on the build machine. Building it for
# $BUILDPLATFORM lets a multi-arch image build without ARM emulation.
FROM --platform=$BUILDPLATFORM golang:alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /tmp/site .
RUN /tmp/site

FROM nginx:alpine

COPY --from=build /src/public /usr/share/nginx/html

EXPOSE 80/tcp
