import os
import torch
from heartlib import HeartMuLa
from config import HEARTMULA_LOCAL_PATH, HEARTMULA_HF_PATH, YOURMT3_LOCAL_PATH, YOURMT3_HF_ID

# Cache for models
_HEARTMULA_MODEL = None
_YOURMT3_MODEL = None
_YOURMT3_PROCESSOR = None

def load_heartmula():
    """Load HeartMuLa model, using local cache if available."""
    global _HEARTMULA_MODEL
    if _HEARTMULA_MODEL is None:
        if os.path.isdir(HEARTMULA_LOCAL_PATH):
            _HEARTMULA_MODEL = HeartMuLa.load(HEARTMULA_LOCAL_PATH)
        else:
            _HEARTMULA_MODEL = HeartMuLa.load(HEARTMULA_HF_PATH)
    return _HEARTMULA_MODEL

def load_yourmt3():
    """Load YourMT3 model and tokenizer, using local cache if available."""
    global _YOURMT3_MODEL, _YOURMT3_PROCESSOR
    if _YOURMT3_MODEL is None:
        try:
            from transformers import AutoModelForSeq2SeqLM, AutoTokenizer
            
            path_to_load = YOURMT3_LOCAL_PATH if os.path.isdir(YOURMT3_LOCAL_PATH) else YOURMT3_HF_ID
            
            _YOURMT3_MODEL = AutoModelForSeq2SeqLM.from_pretrained(path_to_load)
            _YOURMT3_PROCESSOR = AutoTokenizer.from_pretrained(path_to_load)
            
            if torch.cuda.is_available():
                _YOURMT3_MODEL = _YOURMT3_MODEL.to("cuda")
                
        except Exception as e:
            print(f"Error loading YourMT3 model: {e}")
            
    return _YOURMT3_MODEL, _YOURMT3_PROCESSOR
