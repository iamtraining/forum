package web

import (
	"time"

	"github.com/gofiber/fiber"
	"github.com/gofiber/session"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
)

var (
	cfg = session.Config{
		Expiration: 30 * time.Minute,
		GCInterval: 30 * time.Minute,
	}
	s = session.New(cfg)
)

type SessionData struct {
	Form     interface{}
	User     entity.User
	LoggedIn bool
}

func GetSessionData(c *fiber.Ctx) SessionData {
	var data SessionData
	st := s.Get(c)
	data.User, data.LoggedIn = st.Get("user").(entity.User)

	st.Set("form", data.Form)
	if data.Form == nil {
		data.Form = map[string]string{}
	}

	return data
}

func (h *Handler) IsLoggedIn(c *fiber.Ctx) {
	store := s.Get(c)
	id, _ := store.Get("user_id").(uuid.UUID)
	user, err := h.store.User(id)
	if err != nil {
		c.Next()
		return
	}

	store.Set("user", user)
	c.Next()
}

func Login(c *fiber.Ctx, userID uuid.UUID) {
	store := s.Get(c)
	store.Set("user_id", userID)
	store.Save()
}

func Logout(c *fiber.Ctx) {
	store := s.Get(c)
	store.Delete("user_id")
	err := store.Save()
	if err != nil {
		panic(err)
	}
	c.ClearCookie()
	c.Send("you are logged out")
	return
}
