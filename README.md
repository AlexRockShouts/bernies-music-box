# Berne's Music Box Monorepo

This monorepo contains a set of projects for AI music generation and orchestration.

## Project Structure

- `runpod-serverless/`: PyTorch/CUDA environment with FastAPI and RunPod handler for HuggingFace and MuLa models.
- `api-gateway/`: Go-based REST API that acts as an orchestration layer, handling task queuing and file management.
- `ui/`: Modern reactive Svelte (with TypeScript) frontend for user interaction.

## Getting Started

### Prerequisites

- Docker
- Go 1.25+
- Node.js & npm/pnpm/yarn
- Python 3.10+
- NVIDIA GPU (for local CUDA testing, or use RunPod)

### Installation

Refer to the README in each subdirectory for specific installation instructions.

### Orchestration

Use the root `Makefile` for common tasks:

```bash
make build    # Build all components
make run      # Run the system locally (requires appropriate environment)
```

## License

MIT
