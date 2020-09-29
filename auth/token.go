package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Token struct {
	Hash   string
	Expire int64
}

type Claims struct {
	UserID uuid.UUID
	jwt.StandardClaims
}

func CreateToken(c *fiber.Ctx, userID uuid.UUID, secret []byte) (Token, error) {
	var t Token

	exp := time.Now().Add(15 * time.Minute)

	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: exp.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenHash, err := token.SignedString(secret)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error() + "1",
		})
		return Token{}, err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "forum-Token",
		Value:    tokenHash,
		Secure:   false,
		HTTPOnly: true,
	})

	t.Expire, t.Hash = exp.Unix(), tokenHash

	return t, nil
}

func ParseToken(c *fiber.Ctx, secret []byte) (uuid.UUID, error) {
	t := c.Cookies("forum-Token")
	if t == "" {
		return [16]byte{}, fmt.Errorf("token field is empty")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(t, claims, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": err.Error() + "2",
			})
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": err.Error() + "3",
			})
			return [16]byte{}, err
		}
		return [16]byte{}, err
	}

	if !token.Valid {
		c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "invalid token",
		})
		return [16]byte{}, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

func RefreshToken(c *fiber.Ctx, secret []byte) {
	id, err := ParseToken(c, secret)

	if err != nil {
		return
	}

	CreateToken(c, id, secret)
}

func DeleteToken(c *fiber.Ctx) {
	c.ClearCookie("forum-Token")
}
