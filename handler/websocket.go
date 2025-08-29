package handler

import (
	"context"
	"log"
	"sync"
	"time"
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

var (
	wsConnections = make([]*websocket.Conn, 0)
	wsMutex       sync.Mutex
)

func convertUTCToISO(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}
	return ts.UTC().Format(time.RFC3339)
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
			"timestamp": convertUTCToISO(msg.Timestamp),
		}
		if err := conn.WriteJSON(out); err != nil {
			return err
		}
	}
	return nil
}

func WebsocketHandler(conn *websocket.Conn) {
	wsMutex.Lock()
	wsConnections = append(wsConnections, conn)
	wsMutex.Unlock()
	defer func() {
		wsMutex.Lock()
		for i, c := range wsConnections {
			if c == conn {
				wsConnections = append(wsConnections[:i], wsConnections[i+1:]...)
				break
			}
		}
		wsMutex.Unlock()
		conn.Close()
	}()

	// send history messages
	if err := sendHistoryMessages(conn); err != nil {
		log.Println("Send history error:", err)
	}

	for {
		var data map[string]interface{}
		err := conn.ReadJSON(&data)
		if err != nil {
			log.Println("WebSocket read error:", err)
			break
		}
		ts := time.Now().UTC()
		msg := Message{
			Username:  tool.ToString(data["username"], "Anonymous"),
			Message:   tool.ToString(data["message"], ""),
			Timestamp: ts,
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = MsgCol.InsertOne(ctx, msg)
		cancel()
		if err != nil {
			log.Println("Mongo insert error:", err)
		}
		out := map[string]interface{}{
			"username":  msg.Username,
			"message":   msg.Message,
			"timestamp": convertUTCToISO(ts),
		}
		log.Printf("[New message] [Message_time:%s] [User:%s]: [Message:%s]", ts, msg.Username, msg.Message)
		wsMutex.Lock()
		for _, c := range wsConnections {
			c.WriteJSON(out)
		}
		wsMutex.Unlock()
	}
}
