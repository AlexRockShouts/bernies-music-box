import os
import secrets
import uuid
import torch
from pathlib import Path
from fastapi import FastAPI, BackgroundTasks, HTTPException, Depends, status
from fastapi.responses import FileResponse
from fastapi.security import HTTPBasic, HTTPBasicCredentials
from pydantic import BaseModel
from heartlib import HeartMuLa  # adjust import based on actual heartlib API
from basic_pitch import predict  # or use MT3 if you prefer

app = FastAPI(title="Rock/Pop Music Generator API")
security = HTTPBasic()

# Configuration from environment
AUTH_USERNAME = os.getenv("AUTH_USERNAME", "admin")
AUTH_PASSWORD = os.getenv("AUTH_PASSWORD", "password")

def authenticate(credentials: HTTPBasicCredentials = Depends(security)):
    is_correct_username = secrets.compare_digest(credentials.username, AUTH_USERNAME)
    is_correct_password = secrets.compare_digest(credentials.password, AUTH_PASSWORD)
    if not (is_correct_username and is_correct_password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Basic"},
        )
    return credentials.username

class GenerateRequest(BaseModel):
    lyrics: str
    style: str = "indie rock, electric guitar, driving drums, male vocal, E minor"
    duration: int = 180  # seconds
    seed: int = -1

# Simple in-memory storage (use SQLite later)
generations = {}

@app.post("/generate")
async def generate(req: GenerateRequest, background_tasks: BackgroundTasks, username: str = Depends(authenticate)):
    task_id = str(uuid.uuid4())
    generations[task_id] = {"status": "processing", "progress": 0}
    
    background_tasks.add_task(run_generation, task_id, req)
    return {"task_id": task_id, "status": "processing"}

async def run_generation(task_id: str, req: GenerateRequest):
    try:
        # Load model (do this once at startup in production)
        # Note: In a real serverless env, you might want to load this outside the handler
        model = HeartMuLa.load("HeartMuLa/HeartMuLa-oss-3B-happy-new-year")
        
        audio_path = Path(f"output/{task_id}.wav")
        audio_path.parent.mkdir(exist_ok=True)
        
        # Generate audio with style tags
        audio = model.generate(
            lyrics=req.lyrics,
            style=req.style,
            duration=req.duration,
            seed=req.seed if req.seed != -1 else None
        )
        audio.save(str(audio_path))
        
        # MIDI extraction (grooves, notes, chords)
        midi_path = audio_path.with_suffix(".mid")
        # basic_pitch returns note events
        # Note: predict() might be blocking, consider running in a thread pool if needed
        model_output, _ = predict(str(audio_path)) 
        
        # TODO: Add MIDI writing logic using model_output
        # For now, we'll just mock the completion
        
        generations[task_id] = {
            "status": "completed",
            "audio_url": f"/download/{task_id}.wav",
            "midi_url": f"/download/{task_id}.mid"
        }
    except Exception as e:
        generations[task_id] = {"status": "failed", "error": str(e)}

@app.get("/task/{task_id}")
def get_task(task_id: str, username: str = Depends(authenticate)):
    return generations.get(task_id, {"status": "not_found"})

@app.get("/download/{filename}")
def download_file(filename: str):
    file_path = Path("output") / filename
    if file_path.exists():
        return FileResponse(file_path)
    raise HTTPException(status_code=404)

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
