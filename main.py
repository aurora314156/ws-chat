from db.mongo_init import check_mongo_connection
from fastapi import FastAPI, WebSocket
from fastapi.responses import HTMLResponse, FileResponse
from fastapi.staticfiles import StaticFiles
from chat.websocket import websocket_endpoint
from pydantic import BaseModel


class StatusResponse(BaseModel):
    status: str
    message: str

app = FastAPI(
    title="Chat Server API",
    description="Simple chat application with FastAPI + WebSocket + MongoDB",
    version="1.0.0",
    redoc_url="/redoc"
)


# 掛載 static 資料夾
app.mount("/static", StaticFiles(directory="static"), name="static")

@app.get("/", response_model=StatusResponse)
async def root():
    return {"status": "ok", "message": "Chat server running 🚀"}


@app.get("/chat")
async def get():
    with open("static/chat.html", "r", encoding="utf-8") as f:
        return HTMLResponse(f.read())


@app.websocket("/ws")
async def ws(websocket: WebSocket):
    await websocket_endpoint(websocket)


@app.on_event("startup")
async def startup_event():
    await check_mongo_connection()