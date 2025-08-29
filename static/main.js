const chat = document.getElementById("chat");

// Convert UTC ISO string to local time format: "YYYY/MM/DD HH:MM:SS"
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

// -----------------------------
// WebSocket realtime message
// -----------------------------
async function initWebSocket() {
    // get backend url from firebase config
    const res = await fetch("/static/config.json");
    const data = await res.json();
    const BACKEND_URL = data.BACKEND_URL;

    const ws = new WebSocket(`wss://${BACKEND_URL}/ws`);

    ws.onmessage = function(event) {
        const data = JSON.parse(event.data);
        displayMessage(data);
    };

    window.sendMessage = function() {
        const username = document.getElementById("username").value || "Anonymous";
        const message = document.getElementById("message").value.trim();
        if(!message) return;

        ws.send(JSON.stringify({ username, message }));
        document.getElementById("message").value = "";
    };
}

initWebSocket();
