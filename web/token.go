package web

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreateToken(c *fiber.Ctx, userID uuid.UUID, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = userID
	exp := time.Now().Add(time.Hour * 24).Unix()
	claims["exp"] = exp

	t, err := token.SignedString(secret)

	if err != nil {
		return "", err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "forum-Token",
		Value:    t,
		Secure:   false,
		HTTPOnly: true,
	})
	return t, nil
}

func Login(c *fiber.Ctx, userID uuid.UUID, secret []byte) (string, error) {
	store := sessions.Get(c)
	store.Set("user_id", userID)

	token, err := CreateToken(c, userID, []byte("SECRET_KEY"))
	if err == nil {
		store.Set("user_token", token)
	}
	store.Save()

	return token, err
}

func DeleteToken(c *fiber.Ctx) {
	c.ClearCookie("forum-Token")
}
