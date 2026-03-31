import json
import zipfile
from pathlib import Path
from .config import OUTPUT_DIR

def get_task_dir(username: str, task_id: str) -> Path:
    """Get the task-specific output directory."""
    path = Path(OUTPUT_DIR) / username / task_id
    path.mkdir(exist_ok=True, parents=True)
    return path

def save_metadata(task_dir: Path, metadata: dict):
    """Save metadata to a JSON file."""
    with open(task_dir / "metadata.json", "w") as f:
        json.dump(metadata, f, indent=2)

def load_metadata(task_dir: Path) -> dict:
    """Load metadata from a JSON file."""
    metadata_path = task_dir / "metadata.json"
    if not metadata_path.exists():
        return {}
    with open(metadata_path, "r") as f:
        return json.load(f)

def create_zip(zip_path: Path, files: list[Path], metadata_path: Path = None):
    """Create a ZIP archive from a list of files."""
    with zipfile.ZipFile(zip_path, 'w') as zipf:
        for f in files:
            if f.exists():
                zipf.write(f, f.name)
        if metadata_path and metadata_path.exists():
            zipf.write(metadata_path, "metadata.json")
