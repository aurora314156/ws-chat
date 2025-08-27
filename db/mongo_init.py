
import os
from motor.motor_asyncio import AsyncIOMotorClient
from dotenv import load_dotenv
from logger import logger

# Prefer reading MongoDB-related environment variables from the cloud environment.
# If not available, fall back to reading from the local .env file.
if os.getenv("ENV") == "LOCAL":
    logger.info("Local environment detected. Loading .env file...")
    load_dotenv()  # load local .env file.


MONGO_URI = os.environ.get("MONGO_URI")
DB_NAME = os.environ.get("DB_NAME", "chat_db")
COLLECTION_NAME = os.environ.get("COLLECTION_NAME", "messages")

logger.info(f"MONGO_URI: {MONGO_URI}, DB_NAME: {DB_NAME}, COLLECTION_NAME: {COLLECTION_NAME}")

if not MONGO_URI or not DB_NAME or not COLLECTION_NAME:
    logger.error("❌ MONGO_ENV parameters is not set! Please check environment variables.")

# ----------------------------
# MongoDB client
# serverSelectionTimeoutMS 防止DB還沒啟動卡住
# ----------------------------
mongo_client = AsyncIOMotorClient(MONGO_URI, serverSelectionTimeoutMS=5000)
db = mongo_client[DB_NAME]
messages_collection = db[COLLECTION_NAME]

# ----------------------------
# MongoDB connection check at startup
# ----------------------------
async def check_mongo_connection():
    try:
        await mongo_client.admin.command("ping")
        logger.info(f"[✅] MongoDB connection successful. Database '{DB_NAME}' is reachable.")
    except Exception as e:
        logger.error(f"[❌] MongoDB connection failed: {e}")


# ----------------------------
# MongoDB init check at startup
# ----------------------------
async def init_db():
    collections = await db.list_collection_names()
    if "messages" not in collections:
        await db.create_collection("messages")
        logger.info("✅ messages collection created")
    else:
        logger.info("ℹ️ messages collection already exists")

