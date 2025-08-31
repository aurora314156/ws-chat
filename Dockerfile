# 第一階段：編譯 Go 應用程式
# 使用包含完整 Go 工具鏈的映像檔
FROM golang:1.25-alpine AS builder

WORKDIR /app

# 設定 GOPROXY，加速 Go 模組下載
ENV GOPROXY=https://goproxy.io,direct

# 複製 go.mod 和 go.sum 檔案，並下載依賴
COPY go.mod go.sum ./
RUN go mod download

# 複製所有原始碼
COPY . .

# 編譯應用程式，並將可執行檔命名為 'app'
RUN go build -o app main.go

# ---
# 第二階段：運行 Go 應用程式
# 使用極小且安全的 alpine 映像檔作為基礎
FROM alpine:latest

WORKDIR /app

# 安裝 SSL/TLS 憑證，以確保應用程式可以進行安全的網路連線
RUN apk add --no-cache ca-certificates

# 從第一階段複製編譯好的 'app' 可執行檔
COPY --from=builder /app/app .

# 設定環境變數和公開的埠號
ENV PORT=8080
EXPOSE 8080

# 執行可執行檔
CMD ["./app"]