package maple

import (
	"log"
	"testing"
)

func TestNew(t *testing.T) {
	app := New()

	app.Hooks().OnStart(func() error {

		log.Println("OnStart")

		return nil
	})

	err := app.Start()
	if err != nil {
		log.Fatal(err)
	}
}
