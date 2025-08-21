#!/bin/bash
set -e

echo "===== 更新系統套件列表 ====="
sudo apt update -y
sudo apt upgrade -y

echo "===== 安裝必要工具 ====="
sudo apt install -y curl wget git software-properties-common apt-transport-https ca-certificates gnupg lsb-release build-essential

echo "===== 安裝 Mise ====="
curl -fsSL https://get.mise.io | bash || true

echo "===== 安裝 Python 3.12 ====="
sudo add-apt-repository ppa:deadsnakes/ppa -y
sudo apt update -y
sudo apt install -y python3.12 python3.12-venv python3.12-dev
# 設定 python3 預設版本
sudo update-alternatives --install /usr/bin/python3 python3 /usr/bin/python3.12 2
sudo update-alternatives --install /usr/bin/python python /usr/bin/python3.12 2

echo "===== 安裝 pip ====="
python3 -m ensurepip --upgrade
python3 -m pip install --upgrade pip

echo "===== 安裝 FastAPI, Uvicorn, Motor, Jinja2 ====="
pip install fastapi uvicorn motor jinja2

echo "===== 啟動 MongoDB Docker 容器 ====="
docker pull mongo:7
if [ "$(docker ps -a -q -f name=mongodb)" ]; then
    docker start mongodb
else
    docker run -d --name mongodb -p 27017:27017 \
      -e MONGO_INITDB_ROOT_USERNAME=admin \
      -e MONGO_INITDB_ROOT_PASSWORD=admin mongo:7
fi

# 等容器啟動
sleep 5

echo "===== 初始化資料庫 ====="
docker exec -i mongodb mongosh -u admin -p admin <<EOF
use chat_db
db.createCollection("messages")
EOF

echo "===== 環境安裝完成 ====="
echo "MongoDB 帳號密碼：admin/admin，資料庫 chat_db"

echo "Python 版本：$(python3 --version)"
echo "Docker 版本：$(docker --version)"


