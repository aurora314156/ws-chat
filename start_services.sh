#!/bin/bash
if [ "$(docker ps -a -q -f name=mongodb)" ]; then
    echo "MongoDB 容器已存在，啟動中..."
    docker start mongodb
else
    echo "MongoDB 容器不存在，創建並啟動中..."
    docker run -d --name mongodb -p 27017:27017 \
      -e MONGO_INITDB_ROOT_USERNAME=admin \
      -e MONGO_INITDB_ROOT_PASSWORD=admin mongo:7

    # 等容器啟動
    sleep 5

    echo "===== 初始化資料庫 ====="
    docker exec -i mongodb mongosh -u admin -p admin <<EOF
use chat_db
db.createCollection("messages")
EOF

    echo "===== 環境安裝完成 ====="
    echo "MongoDB 帳號密碼：admin/admin，資料庫 chat_db"
fi


echo "===== 啟動 FastAPI ====="
uvicorn main:app --reload --host 0.0.0.0 --port 8000

