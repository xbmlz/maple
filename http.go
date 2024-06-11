package maple

import (
	"github.com/gin-contrib/cors"
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

	router := initRouter(app)

	// configure cors
	router.Use(cors.New(cors.Config{
		AllowOrigins: config.AllowedOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

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

	// print banner

	return server.ListenAndServe()
}

func initRouter(app *App) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	return router
}
