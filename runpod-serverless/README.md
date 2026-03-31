# RunPod AI Music Generator (Serverless)

This service provides an AI-powered music generation pipeline using `HeartMuLa` (for audio generation) and `Basic-Pitch` (for audio-to-MIDI transcription). It is designed to run on RunPod's serverless infrastructure but also includes a FastAPI wrapper for local development and testing.

## Features

- **Lyrics-to-Music**: Generate full-length audio tracks based on text prompts and style tags.
- **Audio-to-MIDI**: Automatic extraction of MIDI data (grooves, notes) from generated audio using Basic Pitch or YourMT3+.
- **Dual Mode**: Runs as a RunPod Serverless worker or a standalone FastAPI application.
- **WebSocket Support**: Real-time progress and status updates for long-running generation tasks.
- **Optimized for RunPod**: Docker image includes pre-downloaded models (HeartMuLa, YourMT3+) to eliminate cold-start latency.

## API Documentation

The service exposes several endpoints for managing generation tasks. For a full OpenAPI 3.0 specification, see [openapi.yaml](./openapi.yaml).

### Endpoints

- `POST /generate`: Start a new generation task.
- `GET /task/{task_id}`: Check the status and progress of a task.
- `GET /generations`: List all generation tasks for the authenticated user.
- `GET /download/{username}/{task_id}/{filename}`: Download generated WAV, MIDI, metadata or ZIP files.
- `WS /ws/{task_id}`: Subscribe to real-time task updates.

### Authentication

The FastAPI endpoints are protected via **Basic Auth**.
- **Username**: Any non-empty string (used for per-user storage). `AUTH_USERNAME` is used as a default for RunPod jobs.
- **Password**: `AUTH_PASSWORD` (must be provided via environment).

You can create a `.env` file in the root directory to manage these variables:
```bash
AUTH_USERNAME=your_username
AUTH_PASSWORD=your_password
```
See [.env.example](./.env.example) for a template.

## RunPod Serverless Deployment

### 1. Build and Push Docker Image

```bash
docker build -t your-registry/music-box-runpod:latest .
docker push your-registry/music-box-runpod:latest
```

*Note: The Dockerfile uses `--platform=linux/amd64`. Ensure you are building for the correct target architecture.*

### 2. Configure RunPod Handler

- **Docker Image**: `your-registry/music-box-runpod:latest`
- **Docker Command**: `python3 main.py --serverless`
- **Environment Variables**:
  - `AUTH_USERNAME`: Your chosen username.
  - `AUTH_PASSWORD`: Your chosen password.

### 3. Input Schema for RunPod Jobs

When calling the serverless endpoint via RunPod's API, use the following JSON structure:

```json
{
  "input": {
    "lyrics": "Your lyrics here",
    "style": "indie rock, electric guitar",
    "duration": 120,
    "seed": 42,
    "stems": true,
    "midi": true,
    "midi_model": "pitch"
  }
}
```

## Local Development

### Prerequisites

- Python 3.10+
- `ffmpeg`
- `libsndfile1`

### Installation

Use the `Makefile` to install dependencies:
```bash
make install
```
Or manually:
1. Install system dependencies: `ffmpeg`, `libsndfile1`.
2. Install Python dependencies:
   ```bash
   pip install -r requirements.txt
   pip install git+https://github.com/HeartMuLa/heartlib.git
   ```

### Running Locally

```bash
make run
```
The server will start at `http://localhost:8000`.

### Maintenance

- **Linting**: `make lint` (syntax check)
- **Cleanup**: `make clean` (remove cache and outputs)

## Model Credits

- **HeartMuLa**: AI music generation model.
- **Basic Pitch**: Spotify's lightweight MIDI transcription model.
