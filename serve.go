package maple

import "github.com/gofiber/fiber/v3"

// ServeConfig defines a configuration struct for apis.Serve().
type ServeConfig struct {

	// HttpAddr is the TCP address to listen for the HTTP server (eg. `127.0.0.1:80`).
	HttpAddr string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

func Serve(a *App, config ServeConfig) error {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	mainAddr := config.HttpAddr

	app := fiber.New(fiber.Config{})

	app.Get("/", func(c fiber.Ctx) error {
		// Send a string response to the client
		return c.SendString("Hello, World 👋!")
	})

	serveEvent := &ServeEvent{
		App:    a,
		Router: app,
	}

	if err := a.OnBeforeServe().Trigger(serveEvent); err != nil {
		return err
	}

	return app.Listen(mainAddr)
}
