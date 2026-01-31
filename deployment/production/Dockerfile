# Stage: builder
FROM golang:1.25.1 AS builder
WORKDIR /app

# Ensure binaries installed by `go install` land in a predictable location
ENV GOBIN=/go/bin
ENV PATH=$GOBIN:$PATH

# Install protobuf compiler for gRPC
# First, try installing unzip. If network fails, continue anyway as we'll handle it differently
RUN apt-get update && apt-get install -y unzip || true

# Download and install protoc
RUN PROTOC_VERSION=28.3 && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    if command -v unzip >/dev/null 2>&1; then \
        unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local bin/protoc && \
        unzip -o protoc-${PROTOC_VERSION}-linux-x86_64.zip -d /usr/local 'include/*'; \
    else \
        python3 -m zipfile -e protoc-${PROTOC_VERSION}-linux-x86_64.zip /tmp/protoc && \
        cp /tmp/protoc/bin/protoc /usr/local/bin/ && \
        cp -r /tmp/protoc/include/* /usr/local/include/; \
    fi && \
    chmod +x /usr/local/bin/protoc && \
    rm -f protoc-${PROTOC_VERSION}-linux-x86_64.zip && \
    rm -rf /var/lib/apt/lists/* /tmp/protoc

# Install Go protobuf plugins for gRPC
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Configure Go environment
RUN go env -w GOPROXY="https://proxy.golang.org,https://goproxy.io,direct" && \
    go env -w GOVCS="*:all"

# Copy dependency files
COPY go.mod go.sum ./

# Generate protobuf files for gRPC (dynamically for all .proto files)
RUN if [ -d proto ] && [ -n "$(find proto -maxdepth 1 -name '*.proto')" ]; then \
        protoc --go_out=. --go_opt=paths=source_relative \
            --go-grpc_out=. --go-grpc_opt=paths=source_relative \
            proto/*.proto && \
        mkdir -p proto/pb && \
        find proto -maxdepth 1 -name '*.pb.go' -exec mv {} proto/pb/ \; && \
        find proto -maxdepth 1 -name '*_grpc.pb.go' -exec mv {} proto/pb/ \;; \
    fi
# Generate Swagger docs before building
RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag init --parseDependency --output ./api-docs/swagger

# Build the application (CGO enabled for Kafka)
RUN go build -o main .

# Clean up sensitive files after build
RUN rm -f .env 
# Build and cache goose CLI in builder stage (static binary for Debian)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go install github.com/pressly/goose/v3/cmd/goose@latest
# Install swag CLI for swagger generation
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go install github.com/swaggo/swag/cmd/swag@latest

# Final stage (unified: contains both migration tools AND application)
# Using Debian-based image for glibc compatibility (required for CGO/Kafka)
FROM debian:bookworm-slim as deploy
WORKDIR /app

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates bash curl make && \
    rm -rf /var/lib/apt/lists/* && \
    update-ca-certificates

# Create migrations directory
RUN mkdir -p ./database/migrations

# Copy migration files from builder
COPY --from=builder /app/database/migrations/ ./database/migrations/

# Remove any non-essential files, keep only .sql and atlas.sum
RUN cd database/migrations && \
    find . -type f ! -name "*.sql" ! -name "atlas.sum" -delete

# Copy Makefile for migration commands
COPY --from=builder /app/Makefile ./Makefile

# Copy goose CLI for migrations
COPY --from=builder /go/bin/goose /usr/local/bin/goose
RUN chmod +x /usr/local/bin/goose

# Copy the main binary (used for both migrations and app)
COPY --from=builder /app/main .

# Copy migration entrypoint script
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh

# Copy Swagger docs
COPY --from=builder /app/api-docs ./api-docs

# Create non-root user for security
RUN useradd -u 1001 -r -g 0 -d /app -s /sbin/nologin \
    -c "Application user" appuser && \
    chown -R appuser:0 /app && \
    chmod -R g=u /app

# Switch to non-root user
USER appuser

# Expose HTTP and gRPC ports
EXPOSE 3013 50063

# Default: run application
# Override with docker-entrypoint.sh for migrations in initContainer
CMD ["/app/main"]
