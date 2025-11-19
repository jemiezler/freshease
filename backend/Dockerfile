# Stage 1: builder
ARG TARGETPLATFORM
FROM --platform=$TARGETPLATFORM golang:1.25-alpine AS builder
WORKDIR /src

# copy module files then download deps for caching
COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on
RUN go mod download

# copy rest of sources
COPY . .

# optional debug
RUN echo "=== go env ===" && go env

# build
RUN set -eux; \
  BUILD_TARGET="."; \
  echo "Building ${BUILD_TARGET}"; \
  CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/backend "${BUILD_TARGET}"

# Stage 2: final minimal image
FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/backend /usr/local/bin/backend
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/backend"]
