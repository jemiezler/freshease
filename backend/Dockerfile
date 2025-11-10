FROM golang:1.25-alpine AS build
WORKDIR /src

# Copy dependency files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

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
