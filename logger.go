package maple

import (
	"log/slog"
)

func (a *App) initLogger() error {
	// TODO
	a.logger = slog.Default()
	return nil
}
