# Cloud Run production Dockerfile
FROM python:3.12-slim

WORKDIR /app

# 安裝依賴
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# 複製程式碼
COPY . .

# Cloud Run 會提供 PORT 環境變數
ENV PORT=8080
EXPOSE 8080

# 啟動 App
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8080"]
