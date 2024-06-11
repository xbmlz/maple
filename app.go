package maple

import (
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
)

// Version of current package
const Version = "v0.0.1"

// App denotes the application.
type App struct {
	mutex sync.Mutex
	// App config
	config Config
	// App logger
	logger *slog.Logger
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
var (
	DefaultDataDir string
	DefaultIsDev   bool
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

	baseDir, isUsingGoRun := inspectRuntime()
	if app.config.DataDir == "" {
		DefaultDataDir = filepath.Join(baseDir, "data")
	}
	if app.config.IsDev == false {
		DefaultIsDev = isUsingGoRun
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

	done := make(chan bool, 1)

	// listen for interrupt signal to gracefully shutdown the application
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh

		done <- true
	}()

	// execute the root command
	go func() {
		// note: leave to the commands to decide whether to print their error
		_ = app.RootCmd.Execute()
		done <- true
	}()

	<-done

	return nil
}

// Hooks returns the hook struct to register hooks.
func (app *App) Hooks() *Hooks {
	return app.hooks
}

// DataDir returns the data directory.
func (app *App) DataDir() string {
	return app.config.DataDir
}

// Logger returns the default app logger.
func (app *App) Logger() *slog.Logger {
	if app.logger == nil {
		return slog.Default()
	}

	return app.logger
}

// inspectRuntime tries to find the base executable directory and how it was run.
func inspectRuntime() (baseDir string, withGoRun bool) {
	if strings.HasPrefix(os.Args[0], os.TempDir()) {
		// probably ran with go run
		withGoRun = true
		baseDir, _ = os.Getwd()
	} else {
		// probably ran with go build
		withGoRun = false
		baseDir = filepath.Dir(os.Args[0])
	}
	return
}
