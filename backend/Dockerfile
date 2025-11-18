FROM golang:1.25-alpine AS builder
WORKDIR /src

# cache modules
COPY go.mod go.sum ./
RUN go mod download

# copy sources
COPY . .

# build for Linux; Buildx will produce arch-specific binaries automatically
# Do not hardcode GOARCH here so buildx can produce arm64/amd64 variants
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /app/backend ./cmd/server

# --- final minimal image ---
FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/backend /usr/local/bin/backend
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/backend"]
