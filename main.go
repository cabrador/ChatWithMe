package main

import (
	"log"
	"os"

	"github.com/petr-hanzl/chatwithme/ai"
	"github.com/petr-hanzl/chatwithme/db"
	"github.com/petr-hanzl/chatwithme/handler"
	"github.com/petr-hanzl/chatwithme/logger"
	"github.com/petr-hanzl/chatwithme/middleware"

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
	chatHandler := handler.NewChatHandler(ai.MakeChatGenerator(db))
	chatGroup.POST("/:personaId", chatHandler.ChatPostHandler)
	chatGroup.GET("", chatHandler.ChatGetHandler)

	log.Fatal(app.Start(":3000"))
}
