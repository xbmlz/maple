package maple

import (
	"github.com/gofiber/fiber/v3"
)

type ServeEvent struct {
	App    *App
	Router *fiber.App
}
