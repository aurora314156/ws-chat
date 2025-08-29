package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"ws-chat/db"
	"ws-chat/handler"
	"ws-chat/logger"
	"ws-chat/wsconn"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var wsManager = wsconn.New()

func main() {
	// init MongoDB
	logger.Info("========== Init Mongo DB Go ==========")
	if err := db.InitMongo(); err != nil {
		logger.Error("MongoDB init failed: %v", err)
	}

	logger.Info("========== Live chat server init starting==========")
	engine := gin.Default()

	engine.Static("/static", "./static") // set up static file serving

	engine.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
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
			logger.Error("WebSocket upgrade error:", err)
			return
		}
		handler.WebsocketHandler(conn, wsManager)
	})

	// Start HTTP server and wait for interrupt signal for graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	// Channel to listen for interrupt or terminate signal
	quit := make(chan struct{})
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server run error:", err)
		}
		close(quit)
	}()

	// Wait here until server exits (in production, you should use os/signal to catch SIGINT/SIGTERM)
	<-quit

	logger.Info("Live chat shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close all WebSocket
	wsManager.CloseAll(ctx)

	// close HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown error:", err)
	}

	logger.Info("Live chat server exited properly")

}
