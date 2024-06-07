package maple

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type HTTPServerConfig struct {

	// HttpAddr is the TCP address to listen for the HTTP server (eg. `127.0.0.1:80`).
	HttpAddr string

	// AllowedOrigins is an optional list of CORS origins (default to "*").
	AllowedOrigins []string
}

func NewHTTPServer(app *App, config HTTPServerConfig) error {
	if len(config.AllowedOrigins) == 0 {
		config.AllowedOrigins = []string{"*"}
	}

	mainAddr := config.HttpAddr

	router := gin.New()

	server := &http.Server{
		ReadTimeout:       10 * time.Minute,
		ReadHeaderTimeout: 30 * time.Second,
		Handler:           router,
		Addr:              mainAddr,
	}

	hook := HTTPServerHook{
		App:    app,
		Router: router,
		Server: server,
	}

	// OnBeforeServer
	app.hooks.executeOnBeforeServer(hook)

	return server.ListenAndServe()
}
