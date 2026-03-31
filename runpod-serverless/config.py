import os
import secrets
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

# Configuration from environment
AUTH_USERNAME = os.getenv("AUTH_USERNAME", "admin")
AUTH_PASSWORD = os.getenv("AUTH_PASSWORD", "password")

# Model paths
HEARTMULA_LOCAL_PATH = "ckpt/HeartMuLa-oss-3B"
HEARTMULA_HF_PATH = "HeartMuLa/HeartMuLa-oss-3B-happy-new-year"

YOURMT3_LOCAL_PATH = "ckpt/yourmt3"
YOURMT3_HF_ID = "mimbres/yourmt3-plus-full"

# Storage paths
OUTPUT_DIR = "output"

def check_credentials(username: str, password: str) -> bool:
    """Verify password against AUTH_PASSWORD. Any non-empty username is accepted."""
    is_correct_password = secrets.compare_digest(password, AUTH_PASSWORD)
    return is_correct_password and len(username) > 0
