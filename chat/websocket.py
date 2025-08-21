from fastapi import WebSocket, WebSocketDisconnect
from motor.motor_asyncio import AsyncIOMotorClient
from datetime import datetime
import asyncio
import logging

# log setting
logger = logging.getLogger("chat") 
logger.setLevel(logging.INFO)
formatter = logging.Formatter("[%(asctime)s] %(levelname)s: %(message)s")

ch = logging.StreamHandler()
ch.setFormatter(formatter)
logger.addHandler(ch)

# MongoDB 連線設定
MONGO_USER = "admin"
MONGO_PASS = "admin"
MONGO_HOST = "localhost"
MONGO_PORT = 27017
DB_NAME = "chat_db"
COLLECTION_NAME = "messages"

# 建立 MongoDB client
mongo_client = AsyncIOMotorClient(
    f"mongodb://{MONGO_USER}:{MONGO_PASS}@{MONGO_HOST}:{MONGO_PORT}"
)
db = mongo_client[DB_NAME]
messages_collection = db[COLLECTION_NAME]

# 保存 WebSocket 連線
connections = []

async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    connections.append(websocket)
    try:
        # 發送歷史訊息
        async for msg in messages_collection.find().sort("timestamp", 1):
            ts = msg.get("timestamp")
            await websocket.send_json({
                "username": msg.get("username", "Anonymous"),
                "message": msg.get("message", ""),
                "timestamp": ts.isoformat() if isinstance(ts, datetime) else ""
            })

        # 接收新訊息
        while True:
            data = await websocket.receive_json()
            timestamp = datetime.utcnow().isoformat()
            logger.info(f"[New message] [Message_time:{timestamp}] [User:{data.get('username')}]: [Message:{data.get('message')}]")
            # 存 MongoDB
            await messages_collection.insert_one({
                "username": data.get("username", "Anonymous"),
                "message": data.get("message", ""),
                "timestamp": timestamp
            })

            # 傳給前端時轉成可序列化字串
            data_to_send = {
                "username": str(data.get("username", "Anonymous")),
                "message": str(data.get("message", "")),
                "timestamp": timestamp
            }
            for conn in connections:
                await conn.send_json(data_to_send)
    except WebSocketDisconnect:
        connections.remove(websocket)
