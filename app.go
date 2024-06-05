package maple

import (
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
)

const Version = "v0.0.1"

type App struct {
	isDev   bool
	dataDir string

	logger *slog.Logger

	RootCmd *cobra.Command
}

type Config struct {

	// DataDir is the directory that contains the database, log files, etc.
	DataDir string
	// IsDev is true if the application is running in development mode.
	IsDev bool
}

// New creates a new Maple named instance.
func New() *App {
	_, isUsingGoRun := inspectRuntime()

	return NewWithConfig(Config{
		IsDev: isUsingGoRun,
	})
}

func NewWithConfig(config Config) *App {
	if config.DataDir == "" {
		baseDir, _ := inspectRuntime()
		config.DataDir = filepath.Join(baseDir, "data")
	}

	app := &App{
		isDev:   config.IsDev,
		dataDir: config.DataDir,
		RootCmd: &cobra.Command{
			Use:     filepath.Base(os.Args[0]),
			Short:   "Maple CLI",
			Version: Version,
			FParseErrWhitelist: cobra.FParseErrWhitelist{
				UnknownFlags: true,
			},
			// no need to provide the default cobra completion command
			CompletionOptions: cobra.CompletionOptions{
				DisableDefaultCmd: true,
			},
		},
	}

	// replace with a colored stderr writer
	app.RootCmd.SetErr(newErrWriter())

	// parse base flags
	// (errors are ignored, since the full flags parsing happens on Execute())
	_ = app.eagerParseFlags(&config)

	// hide the default help command (allow only `--help` flag)
	app.RootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	return app
}

// Start starts the application, aka. registers the default system
func (a *App) Start() error {
	// register system commands
	a.RootCmd.AddCommand(NewServeCommand(a))
	return a.Execute()
}

func (a *App) Execute() error {

	if err := a.Bootstrap(); err != nil {
		return err
	}

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
		_ = a.RootCmd.Execute()

		done <- true
	}()

	<-done

	return nil
}

func (a *App) DataDir() string {
	return a.dataDir
}

func (a *App) Bootstrap() error {
	// ensure that data dir exist
	if err := os.MkdirAll(a.DataDir(), os.ModePerm); err != nil {
		return err
	}

	// initialize the logger
	if err := a.initLogger(); err != nil {
		return err
	}

	return nil
}

func (a *App) Logger() *slog.Logger {
	if a.logger == nil {
		return slog.Default()
	}
	return a.logger
}

func (a *App) initLogger() error {
	// TODO
	a.logger = slog.Default()
	return nil
}

// eagerParseFlags parses the global app flags before calling pb.RootCmd.Execute().
func (a *App) eagerParseFlags(config *Config) error {
	a.RootCmd.PersistentFlags().StringVar(
		&a.dataDir,
		"dir",
		config.DataDir,
		"the Walle data directory",
	)

	a.RootCmd.PersistentFlags().BoolVar(
		&a.isDev,
		"dev",
		config.IsDev,
		"enable dev mode, aka. printing logs and sql statements to the console",
	)

	return a.RootCmd.ParseFlags(os.Args[1:])
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

// newErrWriter returns a red colored stderr writter.
func newErrWriter() *coloredWriter {
	return &coloredWriter{
		w: os.Stderr,
		c: color.New(color.FgRed),
	}
}

// coloredWriter is a small wrapper struct to construct a [color.Color] writter.
type coloredWriter struct {
	w io.Writer
	c *color.Color
}

// Write writes the p bytes using the colored writer.
func (colored *coloredWriter) Write(p []byte) (n int, err error) {
	colored.c.SetWriter(colored.w)
	defer colored.c.UnsetWriter(colored.w)

	return colored.c.Print(string(p))
}

// newOutWriter returns a green colored stdout writter.
