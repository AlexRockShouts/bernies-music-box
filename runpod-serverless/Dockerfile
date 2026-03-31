# Use PyTorch with CUDA as base image
FROM runpod/pytorch:2.2.1-py3.10-cuda12.1.1-devel-ubuntu22.04

# Set working directory
WORKDIR /app

# === Dependencies (your improved block from before) ===
RUN pip install --no-cache-dir --upgrade pip && \
    pip install --no-cache-dir \
        numpy scipy \
        torch>=2.4 torchaudio>=2.4 \
        transformers>=4.57.0 accelerate>=1.10.1 \
        "huggingface-hub[hf_xet]" hf-transfer \
        sentencepiece protobuf && \
    apt-get update && apt-get install -y --no-install-recommends git-lfs && \
    git lfs install && \
    rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt && \
    pip install --no-cache-dir --no-deps git+https://github.com/HeartMuLa/heartlib.git && \
    pip install --no-cache-dir librosa>=0.11.0

# === Model Downloads ===
RUN mkdir -p ckpt && \
    python -c '
import os
from huggingface_hub import snapshot_download
os.makedirs("ckpt", exist_ok=True)
downloads = [
    ("HeartMuLa/HeartMuLaGen", "ckpt"),
    ("HeartMuLa/HeartMuLa-oss-3B-happy-new-year", "ckpt/HeartMuLa-oss-3B"),
    ("HeartMuLa/HeartCodec-oss", "ckpt/HeartCodec-oss"),
    ("HeartMuLa/HeartTranscriptor-oss", "ckpt/HeartTranscriptor-oss"),
    ("mimbres/yourmt3-plus-full", "ckpt/yourmt3"),
]
for repo_id, local_dir in downloads:
    print(f"Downloading {repo_id} → {local_dir}")
    snapshot_download(repo_id=repo_id, local_dir=local_dir, allow_patterns=["*"])
    print(f"✅ Finished {repo_id}\n")
print("All models downloaded.")
'

# Create output directory
RUN mkdir -p output

# Expose FastAPI port
EXPOSE 8000

# Run the application
# Use --serverless if launched on RunPod, or default to FastAPI for local dev
CMD ["python3", "main.py", "--serverless"]
