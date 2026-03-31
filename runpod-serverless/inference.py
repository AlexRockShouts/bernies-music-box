import torch
import anyio
import uuid
from pathlib import Path
from basic_pitch.inference import predict
from basic_pitch import ICASSP_2022_MODEL_PATH
from models import load_heartmula, load_yourmt3
from storage import get_task_dir, save_metadata, create_zip

# Shared task state
tasks: dict[str, dict] = {}

async def update_task_status(task_id: str, data: dict, broadcast_func=None):
    """Update task status and broadcast via WebSocket if available."""
    if task_id in tasks:
        tasks[task_id].update(data)
        if broadcast_func:
            await broadcast_func(task_id, tasks[task_id])

async def run_yourmt3_inference(input_wav: str, output_mid: str):
    """Multi-track music transcription using YourMT3+ with Basic-Pitch fallback."""
    try:
        import librosa
        import numpy as np
        from pretty_midi import PrettyMIDI
        
        model, tokenizer = await anyio.to_thread.run_sync(load_yourmt3)
        if model is None or tokenizer is None:
            raise ValueError("YourMT3 model or tokenizer not loaded")
        
        audio, sr = await anyio.to_thread.run_sync(librosa.load, input_wav, sr=16000, mono=True)
        
        def extract_features(audio, sr):
            hop_length, n_fft, n_mels = 512, 2048, 128
            mel_spec = librosa.feature.melspectrogram(y=audio, sr=sr, n_fft=n_fft, hop_length=hop_length, n_mels=n_mels)
            log_mel_spec = librosa.power_to_db(mel_spec, ref=np.max)
            norm_features = (log_mel_spec - log_mel_spec.min()) / (log_mel_spec.max() - log_mel_spec.min() + 1e-8)
            return torch.tensor(norm_features).T.unsqueeze(0)

        features = await anyio.to_thread.run_sync(extract_features, audio, sr)
        device = next(model.parameters()).device
        inputs = features.to(device)
        
        def generate_tokens(model, inputs):
            outputs = model.generate(inputs_embeds=inputs, max_length=1024, num_beams=1)
            return outputs[0]

        outputs = await anyio.to_thread.run_sync(generate_tokens, model, inputs)
        result_text = tokenizer.decode(outputs, skip_special_tokens=True)
        
        def tokens_to_midi(tokens_text, output_path):
            pm = PrettyMIDI()
            instruments = {}
            current_time = 0.0
            for token in tokens_text.split():
                if token.startswith("time:"):
                    current_time = float(token.split(":")[1])
                elif token.startswith("program:"):
                    prog = int(token.split(":")[1])
                    if prog not in instruments:
                        inst = PrettyMIDI.Instrument(program=prog)
                        pm.instruments.append(inst)
                        instruments[prog] = inst
                elif token.startswith("note:"):
                    try:
                        parts = token.split(":")[1].split(",")
                        pitch, velocity, duration = int(parts[0]), int(parts[1]) if len(parts) > 1 else 100, float(parts[2]) if len(parts) > 2 else 0.5
                        current_inst = list(instruments.values())[-1] if instruments else None
                        if not current_inst:
                            current_inst = PrettyMIDI.Instrument(program=0)
                            pm.instruments.append(current_inst)
                            instruments[0] = current_inst
                        current_inst.notes.append(PrettyMIDI.Note(velocity=velocity, pitch=pitch, start=current_time, end=current_time + duration))
                    except (ValueError, IndexError): continue
            pm.write(output_path)
            return pm

        await anyio.to_thread.run_sync(tokens_to_midi, result_text, output_mid)
            
    except Exception as e:
        print(f"YourMT3+ failed: {e}. Falling back to basic-pitch.")
        _, midi_data, _ = await anyio.to_thread.run_sync(predict, input_wav, ICASSP_2022_MODEL_PATH)
        await anyio.to_thread.run_sync(midi_data.write, output_mid)

async def run_generation(task_id: str, req_params: dict, username: str, broadcast_func=None):
    """Core generation logic: Audio -> Stems -> MIDI -> ZIP."""
    try:
        model = await anyio.to_thread.run_sync(load_heartmula)
        task_dir = get_task_dir(username, task_id)
        audio_path = task_dir / f"{task_id}.wav"
        
        await update_task_status(task_id, {"status": "generating_audio", "progress": 10}, broadcast_func)
        
        metadata = {
            "task_id": task_id, "username": username, "request": req_params,
            "created_at": str(uuid.uuid1().time), "status": "processing"
        }
        await anyio.to_thread.run_sync(save_metadata, task_dir, metadata)
        
        generate_kwargs = {
            "lyrics": req_params["lyrics"], "style": req_params["style"],
            "duration": req_params["duration"], "seed": req_params["seed"] if req_params["seed"] != -1 else None,
            "stems": req_params["stems"]
        }
        
        audio_bundle = await anyio.to_thread.run_sync(model.generate, **generate_kwargs)
        
        stems_paths = []
        if req_params["stems"] and hasattr(audio_bundle, 'stems'):
            await anyio.to_thread.run_sync(audio_bundle.save, str(audio_path.parent / task_id))
            for name in ["vocals", "drums", "bass", "other"]:
                p = audio_path.parent / f"{task_id}_{name}.wav"
                if p.exists(): stems_paths.append(p)
        else:
            await anyio.to_thread.run_sync(audio_bundle.save, str(audio_path))

        generated_files = [audio_path] + stems_paths
        
        if req_params["stems"] and req_params["midi"]:
            await update_task_status(task_id, {"status": "extracting_midi", "progress": 60}, broadcast_func)
            for i, p in enumerate(stems_paths):
                m_path = p.with_suffix(".mid")
                if req_params["midi_model"] == "mt3":
                    await run_yourmt3_inference(str(p), str(m_path))
                else:
                    _, midi_data, _ = await anyio.to_thread.run_sync(predict, str(p), ICASSP_2022_MODEL_PATH)
                    await anyio.to_thread.run_sync(midi_data.write, str(m_path))
                generated_files.append(m_path)
                progress = 60 + int((i + 1) / len(stems_paths) * 30)
                await update_task_status(task_id, {"progress": progress}, broadcast_func)

        zip_path = task_dir / f"{task_id}.zip"
        await anyio.to_thread.run_sync(create_zip, zip_path, generated_files, task_dir / "metadata.json")
        
        metadata["status"] = "completed"
        metadata["files"] = [f.name for f in generated_files]
        await anyio.to_thread.run_sync(save_metadata, task_dir, metadata)

        status_update = {
            "status": "completed", "progress": 100,
            "audio_url": f"/download/{username}/{task_id}/{task_id}.wav",
            "zip_url": f"/download/{username}/{task_id}/{task_id}.zip"
        }
        await update_task_status(task_id, status_update, broadcast_func)
    except Exception as e:
        await update_task_status(task_id, {"status": "failed", "error": str(e)}, broadcast_func)
