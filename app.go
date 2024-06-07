package maple

import (
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
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
	// Command
	RootCmd *cobra.Command
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
		RootCmd: &cobra.Command{
			Use:                filepath.Base(os.Args[0]),
			Short:              "Maple CLI",
			Version:            Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
			CompletionOptions:  cobra.CompletionOptions{DisableDefaultCmd: true},
		},
	}

	// Define hooks
	app.hooks = newHooks(app)

	// Override config if provided
	if len(config) > 0 {
		app.config = config[0]
	}

	app.RootCmd.PersistentFlags().BoolVar(&app.config.IsDev, "dev", DefaultIsDev, "enable dev mode, aka. printing logs and sql statements to the console")
	app.RootCmd.PersistentFlags().StringVar(&app.config.DataDir, "data-dir", DefaultDataDir, "the directory that contains the database, log files, etc.")

	_ = app.RootCmd.ParseFlags(os.Args[1:])
	// Return app
	return app
}

// Start starts the application.
func (app *App) Start() error {
	app.RootCmd.AddCommand(httpServerCommand(app))
	return nil
}

// Hooks returns the hook struct to register hooks.
func (app *App) Hooks() *Hooks {
	return app.hooks
}
