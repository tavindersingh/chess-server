package main

import (
	"tavinder/chess-server/internal/app/auth"
	"tavinder/chess-server/internal/app/game"
	"tavinder/chess-server/internal/app/match"
	"tavinder/chess-server/internal/app/timer"
	"tavinder/chess-server/internal/app/user"
	"time"

	ws "tavinder/chess-server/internal/websocket"

	"github.com/charmbracelet/log"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	log.SetLevel(log.DebugLevel)

	repository := timer.NewInMemoryTimerRepository()
	gameRepository := game.NewCacheGameRepository()
	timerManager := timer.NewTimerManager(repository)
	gameManager := game.NewGameManager(timerManager, gameRepository)
	matchManager := match.NewMatchManager(gameManager)
	wsServer := ws.NewWebSocketServer(matchManager, gameManager)

	app := fiber.New()

	app.Use("/ws", func(c *fiber.Ctx) error {
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

	handleRoutes(app)

	log.Info("Server is running on port 8000")
	app.Listen(":8000")
}

func handleRoutes(app *fiber.App) {
	var jwtManager = auth.NewJwtManager("secret", time.Hour*24*365)
	var userRepository = user.NewInMemoryUserRepository()
	var authService = auth.NewAuthService(jwtManager, userRepository)
	var authHandler = auth.NewAuthHandler(authService)

	var authMiddleware = auth.NewAuthMiddleware(jwtManager)

	var userHandler = user.NewUserHandler(userRepository)

	app.Post("/auth/login/anonymous", authHandler.AnonymousLogin)

	usersGroup := app.Group("/users", authMiddleware.RequireAuth)

	usersGroup.Get("/me", userHandler.CurrentUser)
}
