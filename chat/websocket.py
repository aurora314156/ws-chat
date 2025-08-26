from fastapi import WebSocket, WebSocketDisconnect
from datetime import datetime
from db.mongo_init import messages_collection
from logger import logger




# Save WebSocket connections
connections = []

async def websocket_endpoint(websocket: WebSocket):
    await websocket.accept()
    connections.append(websocket)
    try:
        # send history messages
        async for msg in messages_collection.find().sort("timestamp", 1):
            ts = msg.get("timestamp")
            await websocket.send_json({
                "username": msg.get("username", "Anonymous"),
                "message": msg.get("message", ""),
                "timestamp": ts.isoformat() if isinstance(ts, datetime) else ""
            })

        # receive message
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

            # Convert to a serializable string before sending to the frontend
            data_to_send = {
                "username": str(data.get("username", "Anonymous")),
                "message": str(data.get("message", "")),
                "timestamp": timestamp
            }
            for conn in connections:
                await conn.send_json(data_to_send)
    except WebSocketDisconnect:
        connections.remove(websocket)
