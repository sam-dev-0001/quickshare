# ==========================================
# STAGE 1: Build the statically-linked Go binary
# ==========================================
FROM golang:1.22-alpine AS builder

# Install certificates and update build dependencies
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /app

# Copy modules manifests first for caching
COPY go.mod ./
RUN go mod download

# Copy the entire codebase
COPY . .

# Compile optimized binary with debug symbols stripped
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o /app/quickshare \
    cmd/server/main.go

# ==========================================
# STAGE 2: Construct the hardened runtime container
# ==========================================
FROM alpine:3.19 AS runner

RUN apk update && apk add --no-cache ca-certificates tzdata

WORKDIR /app

# Create a secure, non-privileged system user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create the data folder and set correct owner permissions
RUN mkdir -p /app/data && chown -R appuser:appgroup /app/data

# Copy built binary frombuilder stage
COPY --from=builder /app/quickshare /app/quickshare

# Copy runtime frontend template structures
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/static /app/static

# Handover permissions to the unprivileged appuser
RUN chown -R appuser:appgroup /app/quickshare /app/templates /app/static

# Switch to the non-root execution user
USER appuser

# Document that the container runs on port 3000
EXPOSE 3000

# Set running environments
ENV PORT=3000
ENV NODE_ENV=production

# Execute the self-contained Go server
ENTRYPOINT ["/app/quickshare"]
