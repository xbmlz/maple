package maple

import (
	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v3/log"
	"net/http"
)

type (
	OnStartHandler = func() error
	OnBeforeServer = func(HTTPServerHook) error
)

type HTTPServerHook struct {
	App    *App
	Router *gin.Engine
	Server *http.Server
}

// Hooks is a struct to use it with App
type Hooks struct {
	// Embed App
	app *App

	// Hooks
	onStart        []OnStartHandler
	onBeforeServer []OnBeforeServer
}

func newHooks(app *App) *Hooks {
	return &Hooks{
		app:            app,
		onStart:        make([]OnStartHandler, 0),
		onBeforeServer: make([]OnBeforeServer, 0),
	}
}

// OnStart is a hook to execute user functions after Start.
func (h *Hooks) OnStart(handler ...OnStartHandler) {
	h.app.mutex.Lock()
	h.onStart = append(h.onStart, handler...)
	h.app.mutex.Unlock()
}

// OnBeforeServer is a hook to execute user functions before the server starts.
func (h *Hooks) OnBeforeServer(handler ...OnBeforeServer) {
	h.app.mutex.Lock()
	h.onBeforeServer = append(h.onBeforeServer, handler...)
	h.app.mutex.Unlock()
}

func (h *Hooks) executeOnStart() {
	for _, v := range h.onStart {
		if err := v(); err != nil {
			log.Errorf("failed to call start hook: %v", err)
		}
	}
}

func (h *Hooks) executeOnBeforeServer(callback HTTPServerHook) {
	for _, v := range h.onBeforeServer {
		if err := v(callback); err != nil {
			log.Errorf("failed to call before server hook: %v", err)
		}
	}
}
