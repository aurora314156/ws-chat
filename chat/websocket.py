from fastapi import WebSocket, WebSocketDisconnect
from datetime import datetime, timezone
from db.mongo_init import messages_collection
from logger import logger


def convert_utc_to_local_time(ts):
    if ts:
        if ts.tzinfo is None:
            ts = ts.replace(tzinfo=timezone.utc)
        iso_ts = ts.isoformat().replace("+00:00", "Z")
    else:
        iso_ts = ""
    return iso_ts

# Save WebSocket connections
connections = []

async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    connections.append(websocket)
    try:
        # send history messages
        async for msg in messages_collection.find().sort("timestamp", 1):
            ts = msg.get("timestamp")
            local_time = convert_utc_to_local_time(ts)
            await websocket.send_json({
                "username": msg.get("username", "Anonymous"),
                "message": msg.get("message", ""),
                "timestamp": local_time
            })

        # receive message
        while True:
            data = await websocket.receive_json()
            timestamp = datetime.now(timezone.utc)
            # insert chat message to db
            await messages_collection.insert_one({
                "username": data.get("username", "Anonymous"),
                "message": data.get("message", ""),
                "timestamp": timestamp
            })

            # Convert to a serializable string before sending to the frontend
            chat_message = {
                "username": str(data.get("username", "Anonymous")),
                "message": str(data.get("message", "")),
                "timestamp": timestamp.isoformat().replace("+00:00", "Z")
            }
            logger.info(f"[New message] [Message_time:{timestamp}] [User:{data.get('username')}]: [Message:{data.get('message')}]")
            for conn in connections:
                await conn.send_json(chat_message)
    except WebSocketDisconnect:
        connections.remove(websocket)


async def show_message_history():
    msgs = []
    async for msg in messages_collection.find().sort("timestamp", 1):
        ts = msg.get("timestamp")
        local_time = convert_utc_to_local_time(ts)
        msgs.append({
            "username": msg.get("username", "Anonymous"),
            "message": msg.get("message", ""),
            "timestamp": local_time
        })
    return msgs