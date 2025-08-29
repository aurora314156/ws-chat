package handler

import (
	"context"
	"time"
	"ws-chat/logger"
	"ws-chat/tool"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var MsgCol *mongo.Collection

type Message struct {
	Username  string    `bson:"username" json:"username"`
	Message   string    `bson:"message" json:"message"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

type WSConnection struct {
	Conn *websocket.Conn
}

func sendHistoryMessages(conn *websocket.Conn) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := MsgCol.Find(ctx, bson.D{}, nil)
	if err != nil {
		return err
	}
	var msgs []Message
	if err := cur.All(ctx, &msgs); err != nil {
		return err
	}
	for _, msg := range msgs {
		out := map[string]interface{}{
			"username":  msg.Username,
			"message":   msg.Message,
			"timestamp": tool.ConvertUTCToISO(msg.Timestamp),
		}
		if err := conn.WriteJSON(out); err != nil {
			return err
		}
	}
	return nil
}

func WebsocketHandler(conn *websocket.Conn, wsManager interface {
	Add(*websocket.Conn)
	Remove(*websocket.Conn)
	Broadcast(map[string]interface{})
}) {
	wsManager.Add(conn)
	defer func() {
		wsManager.Remove(conn)
		conn.Close()
	}()

	// send history messages
	if err := sendHistoryMessages(conn); err != nil {
		logger.Error("Send history error:", err)
	}

	// send and receive messages
	for {
		// read message
		var data map[string]interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			logger.Error("WebSocket read error:", err)
			break
		}
		ts := time.Now().UTC()
		// store message to MongoDB
		msg := Message{
			Username:  tool.ToString(data["username"], "Anonymous"),
			Message:   tool.ToString(data["message"], ""),
			Timestamp: ts,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = MsgCol.InsertOne(ctx, msg)
		cancel()
		if err != nil {
			logger.Error("Mongo insert error:", err)
		}
		// broadcast to all ws connections
		out := map[string]interface{}{
			"username":  msg.Username,
			"message":   msg.Message,
			"timestamp": tool.ConvertUTCToISO(ts),
		}
		logger.Info("[New message] [Message_time:%s] [User:%s]: [Message:%s]", ts, msg.Username, msg.Message)
		wsManager.Broadcast(out)
	}
}
