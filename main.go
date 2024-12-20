package main

import (
	"log"
	"os"

	"github.com/petr-hanzl/chatwithme/internal/db"
	"github.com/petr-hanzl/chatwithme/internal/handler"
	"github.com/petr-hanzl/chatwithme/internal/logger"
	"github.com/petr-hanzl/chatwithme/internal/middleware"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	os.Getenv("OPENAI_API_TOKEN")

	database, err := db.MakeDb()
	if err != nil {
		panic(err)
	}

	app := echo.New()
	app.Use(middleware.Logger(logger.NewLogger("DEBUG", "Logger Middleware")))

	chatHandler := handler.NewChatHandler(db.MakeChatController(database))

	// API
	apiGroup := app.Group("/api/v1")
	chatGroup := apiGroup.Group("/chat")
	chatGroup.POST("/:personaId", chatHandler.PersonaPostHandler)

	log.Fatal(app.Start(":4000"))
}
