package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"ws-chat/db"
	"ws-chat/handler"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func main() {
	// init MongoDB
	log.Println("========== Init Mongo DB Go ==========")
	if err := db.InitMongo(); err != nil {
		log.Fatalf("MongoDB init failed: %v", err)
	}

	log.Println("========== Live chat server init starting==========")

	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/chat")
	})

	engine.GET("/chat", func(c *gin.Context) {
		content, err := ioutil.ReadFile("static/index.html")
		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to load index.html")
			return
		}
		c.Data(http.StatusOK, "text/html; charset=utf-8", content)
	})

	engine.GET("/ws", func(c *gin.Context) {
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		}
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Println("WebSocket upgrade error:", err)
			return
		}
		go handler.WebsocketHandler(conn)
	})

	engine.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, StatusResponse{
			Status:  "ok",
			Message: "Live Chat Server Running",
		})
	})

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	log.Println("live-chat started on :8080")
	if err := engine.Run(":8080"); err != nil {
		log.Fatalf("live-chat server error: %v", err)
	}
}
