FROM golang:1.25-alpine AS build

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /src

# Copy go.mod and go.sum first for better Docker layer caching
COPY go.mod go.sum ./

# Download all dependencies (this downloads entgo.io/ent which includes the CLI tool)
RUN go mod download

# Verify ent package is available
RUN go list -m entgo.io/ent

# Copy all source code
COPY . .

# Generate Ent code from schema
# The generate.go file contains: //go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
RUN go generate ./ent

# Verify Ent code was generated (check for ent package)
RUN go list ./ent || (echo "ERROR: ent package not found after generation" && ls -la ent/ && exit 1)

# Build the application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./main.go

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Copy the binary from build stage
COPY --from=build /src/server /app/server

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/server"]
