# LearnX Pedagogy Microservice Dockerfile
FROM python:3.11-slim

# Set working directory
WORKDIR /app

# Set environment variables
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PYTHONPATH=/app

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements and install
COPY person_c/requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir -r requirements.txt

# Copy application files
COPY person_c /app/person_c
COPY domain_knowledge_base.json /app/domain_knowledge_base.json
COPY domain_knowledge_base.md /app/domain_knowledge_base.md
COPY .env.example /app/.env.example

# Expose microservice port (default: 8001)
EXPOSE 8001

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8001/health || exit 1

# Start FastAPI server via Uvicorn
CMD ["uvicorn", "person_c.api:app", "--host", "0.0.0.0", "--port", "8001"]
