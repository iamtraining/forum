package web

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/auth"
	"github.com/iamtraining/forum/entity"
)

type SessionData struct {
	Form     interface{}
	User     entity.User
	LoggedIn bool
}

/*
func (h *Handler) User(c *fiber.Ctx) {
	st := s.Get(c)
	userID := st.Get("user_id")
	if userID == nil {
		return
	}

	id, err := uuid.Parse(fmt.Sprintf("%s", userID))
	if err != nil {
		return
	}

	user, err := h.store.User(id)

	c.Locals("user", user)

}
*/

func (h *Handler) GetSessionData(c *fiber.Ctx) SessionData {
	var data SessionData

	t := c.Locals("user").(*jwt.Token)
	claims := t.Claims.(jwt.MapClaims)
	id := claims["UserID"].(uuid.UUID)
	fmt.Println(id)
	user, err := h.store.User(id)
	if err != nil {
		data.LoggedIn = false
	} else {
		data.LoggedIn = true
		data.User = user
	}

	if data.Form == nil {
		data.Form = map[string]string{}
	}

	fmt.Println(data.User.Username)

	return data
}

/*
func Auth(c *fiber.Ctx) error {
	if !IsLoggedIn(c) {
		c.Redirect("/login")
	}
	return c.Next()
}

func IsLoggedIn(c *fiber.Ctx) bool {
	store := s.Get(c)
	userID := store.Get("user_id")
	token := store.Get("user_token")
	if userID == nil {
		t := c.Cookies("forum-Token")
		if t == "" {
			c.Cookie(&fiber.Cookie{
				Name:     "forum-Token",
				Value:    fmt.Sprintf("%s", token),
				Secure:   false,
				HTTPOnly: true,
			})
		}
		return true
	}
	return false
}
*/

func (h *Handler) Auth(c *fiber.Ctx) error {
	id := c.Locals("user_id").(uuid.UUID)
	fmt.Println(id)
	user, err := h.store.User(id)
	if err != nil {
		c.Next()
	}

	c.Locals("user", user)
	return c.Next()
}

func Login(c *fiber.Ctx, userID uuid.UUID, secret []byte) (auth.Token, error) {
	c.Locals("user_id", userID)
	token, err := auth.CreateToken(c, userID, secret)
	if err == nil {
		c.Locals("user_token", token.Hash)
		c.Locals("token_expiry", token.Expire)
	}

	return token, err
}

func Logout(c *fiber.Ctx) {
	c.ClearCookie()
	return
}
