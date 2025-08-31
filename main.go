package main

import (
	"context"
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
	msgCol, err := db.InitMongo()
	if err != nil || msgCol == nil {
		logger.Error("MongoDB init failed or collection is nil: %v", err)
	}

	logger.Info("========== Live chat server init starting==========")
	engine := gin.Default()

	// Serve static files
	engine.Static("/static", "./static") // set up static file serving
	
	engine.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, StatusResponse{
			Status:  "success",
			Message: "pong",
		})
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
		handler.WebsocketHandler(conn, wsManager, msgCol)
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
