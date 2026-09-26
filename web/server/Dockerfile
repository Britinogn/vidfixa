# Build stage
FROM golang:1.27-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -tags netgo -ldflags="-s -w" \
    -o app ./cmd/server


# Runtime stage
FROM python:3.13-slim-bookworm

WORKDIR /app

# Install FFmpeg
RUN apt-get update \
    && apt-get install -y --no-install-recommends ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Install yt-dlp
RUN pip install --no-cache-dir yt-dlp

# Copy Go application
COPY --from=builder /app/app .

EXPOSE 10000

CMD ["./app"]