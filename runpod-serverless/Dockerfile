FROM --platform=linux/amd64 runpod/pytorch:2.2.1-py3.10-cuda12.1.1-devel-ubuntu22.04

# Install system dependencies + fast download tools
RUN apt-get update && apt-get install -y --no-install-recommends \
        git-lfs \
    && git lfs install \
    && rm -rf /var/lib/apt/lists/*

# Install Python packages (compatible with Python 3.10)
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
        protobuf

# Copy and install project requirements
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir --no-deps git+https://github.com/HeartMuLa/heartlib.git && \
    pip install --no-cache-dir librosa>=0.11.0

# Download all models
RUN mkdir -p ckpt
RUN HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt HeartMuLa/HeartMuLaGen
RUN HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartMuLa-oss-3B HeartMuLa/HeartMuLa-oss-3B-happy-new-year
RUN HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartCodec-oss HeartMuLa/HeartCodec-oss
RUN HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/HeartTranscriptor-oss HeartMuLa/HeartTranscriptor-oss
RUN HF_HUB_ENABLE_HF_TRANSFER=1 huggingface-cli download --local-dir ckpt/yourmt3 mimbres/yourmt3-plus-full

# Create output directory
RUN mkdir -p output

# Expose FastAPI port
EXPOSE 8000

# Run the application
# Use --serverless if launched on RunPod, or default to FastAPI for local dev
CMD ["python3", "main.py", "--serverless"]
