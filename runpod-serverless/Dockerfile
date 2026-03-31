# Use PyTorch with CUDA as base image
FROM runpod/pytorch:2.2.1-py3.10-cuda12.1.1-devel-ubuntu22.04

# Set working directory
WORKDIR /app

# Update and install system dependencies
# Added wget to help with any manual downloads if needed
RUN apt-get update && apt-get install -y \
    git \
    ffmpeg \
    libsndfile1 \
    git-lfs \
    wget \
    && rm -rf /var/lib/apt/lists/*

RUN git lfs install

# Pre-install heavy dependencies to leverage Docker layer caching better
# Note: large downloads may still take time
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir numpy scipy torch>=2.4 torchaudio>=2.4 transformers>=4.57.0 accelerate>=1.10.1 "huggingface-hub[hf_xet]" hf-transfer sentencepiece protobuf

# Copy requirements for remaining dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir --no-deps git+https://github.com/HeartMuLa/heartlib.git && \
    pip install --no-cache-dir librosa>=0.11.0

# Copy application code and environment template
COPY main.py .
COPY .env.example .

# Pre-download the main models using HF CLI as recommended
# This reduces cold-start latency and ensures all weights are in the image
RUN mkdir -p ckpt && \
    HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt HeartMuLa/HeartMuLaGen && \
    HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartMuLa-oss-3B HeartMuLa/HeartMuLa-oss-3B-happy-new-year && \
    HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartCodec-oss HeartMuLa/HeartCodec-oss && \
    HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartTranscriptor-oss HeartMuLa/HeartTranscriptor-oss && \
    HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/yourmt3 mimbres/yourmt3-plus-full

# Create output directory
RUN mkdir -p output

# Expose FastAPI port
EXPOSE 8000

# Run the application
# Use --serverless if launched on RunPod, or default to FastAPI for local dev
CMD ["python3", "main.py", "--serverless"]
