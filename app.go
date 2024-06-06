package maple

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"sync"
)

// Version of current package
const Version = "v0.0.1"

// App denotes the application.
type App struct {
	mutex sync.Mutex
	// App config
	config Config
	// Hooks
	hooks *Hooks
	// fiber server
	server *fiber.App
}

// Config is a struct holding the application settings.
type Config struct {
	// DataDir is the directory that contains the database, log files, etc.
	DataDir string
	// IsDev is true if the application is running in development mode.
	IsDev bool
}

// Default Config values
const (
	DefaultDataDir = "./data"
	DefaultIsDev   = true
)

// New creates a new Maple named instance.
func New(config ...Config) *App {
	// Create a new app
	app := &App{
		config: Config{},
	}

	// Define hooks
	app.hooks = newHooks(app)

	// Override config if provided
	if len(config) > 0 {
		app.config = config[0]
	}

	// Override default values
	if app.config.DataDir == "" {
		app.config.DataDir = DefaultDataDir
	}

	// Return app
	return app
}

// Start starts the application.
func (app *App) Start() error {
	app.hooks.executeOnStart()
	log.Debug("Started Maple...")
	return nil
}

// Hooks returns the hook struct to register hooks.
func (app *App) Hooks() *Hooks {
	return app.hooks
}
