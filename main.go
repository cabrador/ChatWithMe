package main

import (
	"log"
	"os"

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

	database, err := db.MakeDb()
	if err != nil {
		panic(err)
	}

	app := echo.New()
	app.Use(middleware.Logger(logger.NewLogger("DEBUG", "Logger Middleware")))

	app.GET("", handler.HomeHandler{}.HomeGetHandler)

	chatHandler := handler.NewChatHandler(db.MakeChatController(database))

	// API
	apiGroup := app.Group("/api/v1")
	chatGroup := apiGroup.Group("/chat")
	chatGroup.POST("/:personaId", chatHandler.PersonaPostHandler)
	chatGroup.GET("", chatHandler.ChatGetHandler)

	// FE
	renderGroup := app.Group("/chat")
	renderGroup.GET("/:personaId", chatHandler.PersonaGetHandler)

	log.Fatal(app.Start(":3000"))
}
