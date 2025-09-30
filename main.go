package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ws-chat/controller"
	"ws-chat/db"
	"ws-chat/handler"
	"ws-chat/logger"
	"ws-chat/wsconn"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/supabase-community/supabase-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type StatusResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

var wsManager *wsconn.WSManager
var msgCol *mongo.Collection
var supaClient *supabase.Client

func init() {
	// init websocket manager
	wsManager = wsconn.New()

	// init mongoDB
	logger.Info("========== Init MongoDB ==========")
	msgCol, err := db.InitMongo()
	if err != nil || msgCol == nil {
		logger.Error("MongoDB init failed or collection is nil: %v", err)
	}

	// init Supabase
	logger.Info("========== Init Supabase ==========")
	supaClient, err := db.NewSupabaseClient()
	if err != nil || supaClient == nil {
		log.Fatalf("Failed to initialize Supabase client: %v", err)
	}

}

func main() {
	logger.Info("========== Live chat server init starting==========")
	// init controllers
	userController := controller.NewUserController(supaClient)
	userHandler := handler.NewUserHandler(userController)

	// init gin engine
	engine := gin.Default()

	engine.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
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

	v1 := engine.Group("/api/v1")
	{
		v1.POST("/signup", userHandler.Signup)
	}

	// Start HTTP server and wait for interrupt signal for graceful shutdown
	srv := &http.Server{
		Addr:    ":8080",
		Handler: engine,
	}

	// Channel to listen for interrupt or terminate signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
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
