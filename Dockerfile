# Stage 1
FROM alpine:latest AS build

WORKDIR /opt/HugoApp

# Copy Hugo config into the container Workdir.
COPY . .

# Stage 2
FROM nginx:1.25-alpine AS create

# Set workdir to the NGINX default dir.
WORKDIR /usr/share/nginx/html

# Copy HTML from previous build into the Workdir.
COPY --from=build /opt/HugoApp/public .

# Expose port 80
EXPOSE 80/tcp
