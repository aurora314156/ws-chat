from db.mongo_init import check_mongo_connection
from fastapi import FastAPI, WebSocket
from fastapi.responses import HTMLResponse, RedirectResponse
from fastapi.staticfiles import StaticFiles
from chat.websocket import websocket_endpoint, show_message_history
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


@app.get("/", response_model=StatusResponse)
async def root():
    return RedirectResponse(url="/chat")


@app.get("/chat")
async def get():
    with open("static/chat.html", "r", encoding="utf-8") as f:
        return HTMLResponse(f.read())


@app.websocket("/ws")
async def ws(websocket: WebSocket):
    await websocket_endpoint(websocket)


@app.get("/messages/history")
async def get_history():
    return await show_message_history()
   

@app.on_event("startup")
async def startup_event():
    await check_mongo_connection()