package maple

import (
	"errors"
	"github.com/spf13/cobra"
	"net/http"
)

func httpServerCommand(app *App) *cobra.Command {
	var allowedOrigins []string
	var httpAddr string

	command := &cobra.Command{
		Use:   "serve [domain(s)]",
		Args:  cobra.ArbitraryArgs,
		Short: "Starts the web server (default to 127.0.0.1:8090 if no domain is specified)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				if httpAddr == "" {
					httpAddr = "0.0.0.0:80"
				}
			} else {
				if httpAddr == "" {
					httpAddr = "127.0.0.1:8090"
				}
			}

			err := NewHTTPServer(app, HTTPServerConfig{
				AllowedOrigins: allowedOrigins,
				HttpAddr:       httpAddr,
			})

			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		},
	}

	return command
}
