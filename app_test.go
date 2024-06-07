package maple

import (
	"github.com/gin-gonic/gin"
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	app := New()

	app.Hooks().OnStart(func() error {
		log.Println("OnStart")
		return nil
	})

	app.Hooks().OnBeforeServer(func(hook HTTPServerHook) error {
		hook.Router.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{"hello": "world"})
		})
		return nil
	})

	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
