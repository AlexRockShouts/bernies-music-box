# Use PyTorch with CUDA as base image
FROM runpod/pytorch:2.2.1-py3.10-cuda12.1.1-devel-ubuntu22.04

# Set working directory
WORKDIR /app

# Update and install system dependencies
# Added wget to help with any manual downloads if needed
# Pre-install heavy dependencies to leverage Docker layer caching better
# Note: large downloads may still take time
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
        numpy \
        scipy \
        torch>=2.4 \
        torchaudio>=2.4 \
        transformers>=4.57.0 \
        accelerate>=1.10.1 \
        "huggingface-hub[hf_xet]" \
        hf-transfer \
        sentencepiece \
        protobuf && \
    # Install git-lfs for safety with large model files
    apt-get update && apt-get install -y --no-install-recommends git-lfs && \
    git lfs install && \
    # Clean up
    rm -rf /var/lib/apt/lists/*

# Copy requirements for remaining dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir --no-deps git+https://github.com/HeartMuLa/heartlib.git && \
    pip install --no-cache-dir librosa>=0.11.0



# Copy application code and environment template
COPY main.py .
COPY .env .

# Pre-download the main models using HF CLI as recommended
# This reduces cold-start latency and ensures all weights are in the image
# Download all required models (using hf_xet for speed)
RUN mkdir -p ckpt && \
    python -c '
import os
from huggingface_hub import snapshot_download

os.makedirs("ckpt", exist_ok=True)

repos = [
    ("HeartMuLa/HeartMuLaGen", "ckpt"),
    ("HeartMuLa/HeartMuLa-oss-3B-happy-new-year", "ckpt/HeartMuLa-oss-3B"),
    ("HeartMuLa/HeartCodec-oss", "ckpt/HeartCodec-oss"),
    ("HeartMuLa/HeartTranscriptor-oss", "ckpt/HeartTranscriptor-oss"),
    ("mimbres/yourmt3-plus-full", "ckpt/yourmt3"),
]

for repo_id, local_dir in repos:
    print(f"Downloading {repo_id} to {local_dir}...")
    snapshot_download(
        repo_id=repo_id,
        local_dir=local_dir,
        allow_patterns=["*"],
        # Optional: force Xet if available
        # force_download=False
    )
    print(f"Finished {repo_id}")
'

# Create output directory
RUN mkdir -p output

# Expose FastAPI port
EXPOSE 8000

# Run the application
# Use --serverless if launched on RunPod, or default to FastAPI for local dev
CMD ["python3", "main.py", "--serverless"]
