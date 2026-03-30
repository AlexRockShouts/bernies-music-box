# Use PyTorch with CUDA as base image
FROM pytorch/pytorch:2.2.1-cuda12.1-cudnn8-runtime

# Set working directory
WORKDIR /app

# Update and install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    ffmpeg \
    && rm -rf /var/lib/apt/lists/*

# Copy requirements first to leverage Docker cache
COPY requirements.txt .

# Install dependencies
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir -r requirements.txt

# Install HeartMuLa core (if not in requirements.txt or needs git clone)
# Since it's in requirements.txt, we might not need to clone it manually,
# but if the requirements.txt refers to a local path or we want the latest:
RUN git clone https://github.com/HeartMuLa/heartlib.git /tmp/heartlib && \
    pip install -e /tmp/heartlib && \
    rm -rf /tmp/heartlib

# Copy application code
COPY main.py .

# Create output directory
RUN mkdir -p output

# Expose FastAPI port
EXPOSE 8000

# Run the application
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
