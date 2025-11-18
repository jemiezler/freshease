# syntax=docker/dockerfile:1
FROM golang:1.25-alpine AS builder
WORKDIR /src

# copy module files then download deps for caching
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
RUN go mod download

# copy rest of sources
COPY . .

# debug: show module and go env
RUN echo "=== backend/go.mod ===" && sed -n '1,120p' go.mod || true
RUN echo "=== go env ===" && go env || true
RUN echo "=== go list ./... (before tidy) ===" && go list ./... 2>&1 | sed -n '1,120p' || true

# attempt tidy to ensure modules are correct
RUN go mod tidy -v
RUN echo "=== go list ./... (after tidy) ===" && go list ./... 2>&1 | sed -n '1,200p' || true

# build
RUN set -eux; \
  BUILD_TARGET="."; \
  echo "Building ${BUILD_TARGET}"; \
  CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/backend "${BUILD_TARGET}"

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/backend /usr/local/bin/backend
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/backend"]
