import sys
from app import start_serverless, start_api

if __name__ == "__main__":
    if len(sys.argv) > 1 and sys.argv[1] == "--serverless":
        start_serverless()
    else:
        start_api()
