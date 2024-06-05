package maple

import (
	"errors"
	"github.com/spf13/cobra"
	"net/http"
)

func NewServeCommand(a *App) *cobra.Command {

	var allowedOrigins []string
	var httpAddr string

	command := &cobra.Command{
		Use:          "serve [domain(s)]",
		Args:         cobra.ArbitraryArgs,
		Short:        "Starts the web server (default to 127.0.0.1:8090 if no domain is specified)",
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			// set default listener addresses if at least one domain is specified
			if len(args) > 0 {
				if httpAddr == "" {
					httpAddr = "0.0.0.0:80"
				}
			} else {
				if httpAddr == "" {
					httpAddr = "127.0.0.1:8090"
				}
			}

			err := Serve(a, ServeConfig{
				HttpAddr:       httpAddr,
				AllowedOrigins: allowedOrigins,
			})

			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}

			return err
		},
	}

	command.PersistentFlags().StringSliceVar(
		&allowedOrigins,
		"origins",
		[]string{"*"},
		"CORS allowed domain origins list",
	)

	command.PersistentFlags().StringVar(
		&httpAddr,
		"http",
		"",
		"TCP address to listen for the HTTP server\n(if domain args are specified - default to 0.0.0.0:80, otherwise - default to 127.0.0.1:8090)",
	)
	return command
}
