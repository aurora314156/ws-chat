const chat = document.getElementById("chat");

function formatLocalTime(utcString) {
    if (!utcString) return "";
    const date = new Date(utcString);
    const pad = (n) => n.toString().padStart(2, "0");
    return `${date.getFullYear()}/${pad(date.getMonth()+1)}/${pad(date.getDate())} ` +
           `${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`;
}

function displayMessage(data) {
    const div = document.createElement("div");
    div.textContent = `[${formatLocalTime(data.timestamp)}] ${data.username}: ${data.message}`;
    chat.appendChild(div);
    chat.scrollTop = chat.scrollHeight;
}

function sendMessage(ws) {
    const username = document.getElementById("username").value || "Anonymous";
    const message = document.getElementById("message").value.trim();
    if (!message) return;

    // 在發送前檢查 WebSocket 狀態
    if (ws.readyState === WebSocket.OPEN) {
        ws.send(JSON.stringify({ username, message }));
    } else {
        console.error("WebSocket is not OPEN. Message not sent.");
        // 可以選擇通知用戶連線斷開
    }
    document.getElementById("message").value = "";
}

// -----------------------------
// WebSocket 實時訊息處理
// -----------------------------

async function initWebSocket() {
    // get backend url from firebase config
    const res = await fetch("/static/config.json");
    const data = await res.json();
    const BACKEND_URL = data.BACKEND_URL;
    const ws = new WebSocket(`wss://${BACKEND_URL}/ws`);

    // Log when the connection is successfully opened
    ws.onopen = function() {
        console.log("WebSocket connection established successfully.");
    };

    // Log any errors that occur during the connection attempt
    ws.onerror = function(error) {
        console.error("WebSocket error:", error);
    };

    // Log when the connection is closed
    ws.onclose = function() {
        console.log("WebSocket connection closed.");
    };
    
    ws.onmessage = function(event) {
        const data = JSON.parse(event.data);
        displayMessage(data);
    };
    
    // event listener for send button
    const sendButton = document.getElementById("sendButton");
    sendButton.addEventListener("click", () => sendMessage(ws));  

    // monitor input box
    const messageInput = document.getElementById("message");
    messageInput.addEventListener("keydown", (event) => {
        if (event.key === "Enter") {
            sendMessage(ws);
        }
    });
}

document.addEventListener("DOMContentLoaded", initWebSocket);