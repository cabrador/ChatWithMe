package main

import (
	"chatwithme/db"
	"chatwithme/handler"
	"chatwithme/logger"
	"chatwithme/middleware"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	os.Getenv("OPENAI_API_TOKEN")

	db, err := db.MakeDb()
	if err != nil {
		panic(err)
	}

	app := echo.New()
	app.Use(middleware.Logger(logger.NewLogger("DEBUG", "Logger Middleware")))
	rootGroup := app.Group("/api/v1")
	chatGroup := rootGroup.Group("/chat")
	chatHandler := handler.NewChatHandler(db)
	chatGroup.POST("", chatHandler.ChatPostHandler)

	log.Fatal(app.Start(":3000"))
}
