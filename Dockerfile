# ========================================================
# Multi-Stage Build: LearnX All-in-One Server
# Go Backend + Python Pedagogy Microservice
# Deployable to Render / Cloud as a Single Web Service
# ========================================================

# Stage 1: Build the Go Backend Binary
FROM golang:1.27.1-alpine AS go-builder
WORKDIR /build

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/learnx-server .

# Stage 2: Final Runtime Image (Python 3.11 + Compilers + Go Binary)
FROM python:3.11-slim-bookworm

WORKDIR /app

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PYTHONPATH=/app/pedagogy-service \
    SERVER_HOST=0.0.0.0 \
    PEDAGOGY_SERVICE_URL=http://127.0.0.1:8001

# Install system dependencies, curl for healthchecks, and native compilers
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    ca-certificates \
    gcc \
    g++ \
    && rm -rf /var/lib/apt/lists/*

# Install Python pedagogy dependencies
COPY pedagogy-service/person_c/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r /app/requirements.txt

# Copy Pedagogy Service source files
COPY pedagogy-service /app/pedagogy-service

# Copy compiled Go backend binary and data assets
COPY --from=go-builder /build/learnx-server /app/learnx-server
COPY backend/data /app/data

# Copy Startup Script
COPY start.sh /app/start.sh
RUN chmod +x /app/start.sh

# Expose Render PORT (default 8080)
EXPOSE 8080

# Run entrypoint script
CMD ["/app/start.sh"]
