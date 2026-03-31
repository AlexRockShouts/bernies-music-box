import uuid
import asyncio
from pathlib import Path
from typing import List, Optional
from fastapi import FastAPI, BackgroundTasks, HTTPException, Depends, status, WebSocket, WebSocketDisconnect
from fastapi.responses import FileResponse
from fastapi.security import HTTPBasic, HTTPBasicCredentials
from pydantic import BaseModel
import runpod

from config import AUTH_USERNAME, check_credentials, OUTPUT_DIR
from inference import run_generation, tasks, update_task_status
from storage import load_metadata

app = FastAPI(title="Rock/Pop Music Generator API")
security = HTTPBasic()

# WebSocket tracking
active_websockets: dict[str, List[WebSocket]] = {}

class GenerateRequest(BaseModel):
    lyrics: str
    style: str = "indie rock, electric guitar, driving drums, male vocal, E minor"
    duration: int = 180
    seed: int = -1
    stems: bool = False
    midi: bool = False
    midi_model: str = "pitch"

def authenticate(credentials: HTTPBasicCredentials = Depends(security)):
    if not check_credentials(credentials.username, credentials.password):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Incorrect username or password",
            headers={"WWW-Authenticate": "Basic"},
        )
    return credentials.username

async def broadcast_status(task_id: str, data: dict):
    """Broadcast status updates to all connected WebSockets for a task."""
    if task_id in active_websockets:
        disconnected = []
        for ws in active_websockets[task_id]:
            try:
                await ws.send_json(data)
            except Exception:
                disconnected.append(ws)
        for ws in disconnected:
            if ws in active_websockets[task_id]:
                active_websockets[task_id].remove(ws)

@app.post("/generate")
async def generate(req: GenerateRequest, background_tasks: BackgroundTasks, username: str = Depends(authenticate)):
    task_id = str(uuid.uuid4())
    tasks[task_id] = {
        "status": "processing",
        "progress": 0,
        "username": username,
        "request_params": req.model_dump()
    }
    background_tasks.add_task(run_generation, task_id, req.model_dump(), username, broadcast_status)
    return {"task_id": task_id, "status": "processing"}

@app.websocket("/ws/{task_id}")
async def websocket_endpoint(websocket: WebSocket, task_id: str):
    await websocket.accept()
    if task_id not in tasks:
        await websocket.send_json({"status": "not_found"})
        await websocket.close()
        return
    
    active_websockets.setdefault(task_id, []).append(websocket)
    await websocket.send_json(tasks[task_id])
    
    try:
        while True:
            await websocket.receive_text()
    except WebSocketDisconnect:
        if task_id in active_websockets:
            active_websockets[task_id].remove(websocket)
            if not active_websockets[task_id]:
                del active_websockets[task_id]

@app.get("/task/{task_id}")
def get_task(task_id: str, username: str = Depends(authenticate)):
    task = tasks.get(task_id)
    if not task:
        return {"status": "not_found"}
    if task.get("username") != username:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
    return task

@app.get("/generations")
def list_generations(username: str = Depends(authenticate)):
    user_dir = Path(OUTPUT_DIR) / username
    if not user_dir.exists():
        return []
    
    songs = []
    for task_dir in user_dir.iterdir():
        if task_dir.is_dir():
            meta = load_metadata(task_dir)
            if meta:
                tid = task_dir.name
                meta["downloads"] = {
                    "audio": f"/download/{username}/{tid}/{tid}.wav",
                    "zip": f"/download/{username}/{tid}/{tid}.zip",
                    "metadata": f"/download/{username}/{tid}/metadata.json"
                }
                songs.append(meta)
    return songs

@app.get("/download/{username}/{task_id}/{filename}")
def download_file(username: str, task_id: str, filename: str, authenticated_user: str = Depends(authenticate)):
    if username != authenticated_user:
        raise HTTPException(status_code=status.HTTP_403_FORBIDDEN, detail="Forbidden")
        
    file_path = Path(OUTPUT_DIR) / username / task_id / filename
    if file_path.exists():
        return FileResponse(file_path)
    raise HTTPException(status_code=status.HTTP_404_NOT_FOUND)

def runpod_handler(job):
    job_input = job["input"]
    task_id = job["id"]
    username = job_input.get("username", AUTH_USERNAME)
    
    try:
        req = GenerateRequest(**job_input)
    except Exception as e:
        return {"error": str(e)}
    
    tasks[task_id] = {
        "status": "processing", "progress": 0, "username": username, "request_params": req.model_dump()
    }
    asyncio.run(run_generation(task_id, req.model_dump(), username))
    return tasks.get(task_id)

def start_serverless():
    runpod.serverless.start({"handler": runpod_handler})

def start_api():
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8000)
