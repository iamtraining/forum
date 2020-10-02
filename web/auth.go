package web

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	fiber "github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
)

type SessionData struct {
	Form     interface{}
	User     entity.User
	LoggedIn bool
}

func (h *Handler) Extract(c *fiber.Ctx) (entity.User, error) {
	t := c.Cookies("forum-Token")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(t, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY"), nil
	})
	if err != nil {
		return entity.User{}, err
	}

	err = claims.Valid()

	if err != nil {
		DeleteToken(c)
		return entity.User{}, err
	}

	id, err := uuid.Parse(claims["id"].(string))
	if err != nil {
		DeleteToken(c)
		return entity.User{}, err
	}

	if id == uuid.Nil {
		DeleteToken(c)
		return entity.User{}, fmt.Errorf("null uuid")
	}

	user, _ := h.store.User(id)
	fmt.Println("extract", user)
	return user, nil
}

/*
func (h *Handler) Restricted(c *fiber.Ctx) error {
	local := c.Locals("user").(*jwt.Token)
	claims := local.Claims.(jwt.MapClaims)
	id := claims["id"].(uuid.UUID)
	user, _ := h.store.User(id)
	fmt.Println("extract", user)
}
*/

func Protect() fiber.Handler {
	return jwtware.New(jwtware.Config{
		//ErrorHandler: err,
		SigningKey: []byte("SECRET_KEY"),
	})
}

func err(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Missing or malformed JWT"})
	} else {
		return c.
			Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"message": "Invalid or expired JWT"})
	}
}

func (h *Handler) GetSessionData(c *fiber.Ctx) SessionData {
	var data SessionData

	var err error

	data.User, err = h.Extract(c)
	if err != nil {
		data.LoggedIn = false
	} else {
		data.LoggedIn = true
	}

	if data.Form == nil {
		data.Form = map[string]string{}
	}

	fmt.Println(data.User.Username)

	return data
}

func Logout(c *fiber.Ctx) {
	c.ClearCookie()
	return
}
