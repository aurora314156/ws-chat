
import os
from motor.motor_asyncio import AsyncIOMotorClient
from dotenv import load_dotenv
from logger import logger

# Prefer reading MongoDB-related environment variables from the cloud environment.
# If not available, fall back to reading from the local .env file.
if not all([os.environ.get("MONGO_URI"), os.environ.get("DB_NAME"), os.environ.get("COLLECTION_NAME")]):
    load_dotenv()  # load local .env file.


MONGO_URI = os.environ.get("MONGO_URI", "mongodb://admin:admin@mongodb:27017/chat_db?authSource=admin")
DB_NAME = os.environ.get("DB_NAME", "chat_db")
COLLECTION_NAME = os.environ.get("COLLECTION_NAME", "messages")

# ----------------------------
# MongoDB client
# ----------------------------
mongo_client = AsyncIOMotorClient(MONGO_URI)
db = mongo_client[DB_NAME]
messages_collection = db[COLLECTION_NAME]

# ----------------------------
# MongoDB connection check at startup
# ----------------------------
async def check_mongo_connection():
    try:
        # 使用 serverSelectionTimeoutMS 防止卡住
        await mongo_client.admin.command("ping")
        logger.info(f"[✅] MongoDB connection successful. Database '{DB_NAME}' is reachable.")
    except Exception as e:
        logger.error(f"[❌] MongoDB connection failed: {e}")
