# syntax=docker/dockerfile:1

####################
# Builder stage
####################
FROM golang:1.25-alpine AS builder
WORKDIR /src

# copy go.mod first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# copy the backend source only (workflow uses context ./backend)
COPY . .

# Build a static linux binary (CGO disabled). Buildx will handle CPU arch variants.
RUN set -eux; \
  # if main.go exists in root, build '.'; otherwise try common cmd/* dirs
  if [ -f "./main.go" ]; then \
    BUILD_TARGET="."; \
  elif [ -d ./cmd ] && [ -f ./cmd/genmodule/main.go ]; then \
    # example fallback; you can add more if needed
    BUILD_TARGET="./cmd/genmodule"; \
  else \
    BUILD_TARGET="."; \
  fi; \
  echo "Building ${BUILD_TARGET}"; \
  CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/backend "${BUILD_TARGET}"

####################
# Final image
####################
FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/backend /usr/local/bin/backend

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/backend"]
