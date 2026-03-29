# Dockerfile for ML FastAPI service
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt ./
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

ENV PYTHONUNBUFFERED=1

# Expose default port; Railway will set PORT
EXPOSE 8000

# Use Railway's PORT if set, otherwise 8000
CMD ["sh", "-c", "uvicorn mental_health_ml.serve:app --host 0.0.0.0 --port ${PORT:-8000}"]
