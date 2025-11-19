# syntax=docker/dockerfile:1.5
# Multi-platform build for Flutter web
ARG TARGETPLATFORM=linux/amd64

# ---------- Build stage ----------
FROM --platform=$TARGETPLATFORM dart:stable AS build
WORKDIR /app

# Install Flutter
RUN git clone https://github.com/flutter/flutter.git /usr/local/flutter -b stable
ENV PATH="/usr/local/flutter/bin:/usr/local/flutter/bin/cache/dart-sdk/bin:$PATH"

# Enable flutter web
RUN flutter channel stable && \
    flutter upgrade && \
    flutter config --enable-web

# Copy pubspec and lockfile first for caching
COPY pubspec.* ./

# Get dependencies
RUN flutter pub get --offline || flutter pub get

# Copy source code
COPY . .

# Build Flutter web
RUN flutter build web --release

# ---------- Serve stage ----------
FROM --platform=$TARGETPLATFORM nginx:stable-alpine AS runtime
WORKDIR /usr/share/nginx/html

# Remove default nginx files
RUN rm -rf ./*

# Copy build output from Flutter
COPY --from=build /app/build/web ./

# Optional: Copy custom nginx config (if you have one)
# COPY nginx.conf /etc/nginx/nginx.conf

# Expose port
EXPOSE 80

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
