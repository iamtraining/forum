package web

import (
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
)

var (
	cfg = session.Config{
		Expiration: 30 * time.Minute,
		GCInterval: 30 * time.Minute,
	}
	s = session.New(cfg)
)

func IsLoggedIn(c *fiber.Ctx) bool {
	store := s.Get(c)
	key := store.Get("islogined")
	if key == nil {
		return false
	} else {
		return true
	}

}

func IsAuthed(c *fiber.Ctx) {
	IsLoggedIn(c)
	c.Next()
}
