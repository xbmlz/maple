package maple

import "github.com/gofiber/fiber/v3/log"

type (
	OnStartHandler = func() error
)

// Hooks is a struct to use it with App
type Hooks struct {
	// Embed App
	app *App

	// Hooks
	onStart []OnStartHandler
}

func newHooks(app *App) *Hooks {
	return &Hooks{
		app:     app,
		onStart: make([]OnStartHandler, 0),
	}
}

// OnStart is a hook to execute user functions after Start.
func (h *Hooks) OnStart(handler ...OnStartHandler) {
	h.app.mutex.Lock()
	h.onStart = append(h.onStart, handler...)
	h.app.mutex.Unlock()
}

func (h *Hooks) executeOnStart() {
	for _, v := range h.onStart {
		if err := v(); err != nil {
			log.Errorf("failed to call start hook: %v", err)
		}
	}
}
