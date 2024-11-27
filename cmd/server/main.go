package main

import (
	"tavinder/chess-server/internal/app/game"
	"tavinder/chess-server/internal/app/match"
	"tavinder/chess-server/internal/app/timer"
	"tavinder/chess-server/internal/events"
	ws "tavinder/chess-server/internal/websocket"

	"github.com/charmbracelet/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	eventBus := events.NewEventBus()

	repository := timer.NewInMemoryTimerRepository()
	timerManager := timer.NewTimerManager(repository, eventBus)
	gameManager := game.NewGameManager(timerManager, eventBus)
	matchManager := match.NewMatchManager(gameManager)
	wsServer := ws.NewWebSocketServer(matchManager)

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/game/:id", websocket.New(wsServer.HandleWebSocketConnection))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	log.Info("Server is running on port 8000")
	app.Listen(":8000")
}
