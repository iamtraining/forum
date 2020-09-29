package auth

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type AuthSettings struct {
	// отвечает за пропуск middleware [nil]
	Filter func(*fiber.Ctx) bool

	// функция которая выполняется при неправильном токене
	// или для ошибки [401]
	Error func(*fiber.Ctx, error)

	// функция для действительного токена [nil]
	Successful func(*fiber.Ctx)

	Key    interface{}
	Keys   map[string]interface{}
	CtxKey string
	Claims jwt.Claims
	Method string

	// извлечение метода из запроса header:Authorization
	// header:name, parameter:name, query:name, cookie:name
	TFR string

	// для заголовка авторизации [Bearer]
	AuthScheme string

	keyFunc jwt.Keyfunc
}

func Authentificate(config ...AuthSettings) fiber.Handler {
	cfg := AuthSettings{}
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.Error == nil {
		cfg.Error = func(c *fiber.Ctx, err error) {
			if err.Error() == "missing or malformed jwt" {
				c.Status(fiber.StatusBadRequest)
				c.SendString(err.Error())
			} else {
				c.Status(fiber.StatusUnauthorized)
				c.SendString("invalid or expired jwt")
			}
		}
	}

	if cfg.Successful == nil {
		cfg.Successful = func(c *fiber.Ctx) {
			c.Next()
		}
	}

	if cfg.Key == nil && len(cfg.Keys) == 0 {
		log.Fatal("jwt requeres signing key")
	}

	if cfg.Method == "" {
		cfg.Method = "HS256"
	}

	if cfg.Claims == nil {
		cfg.Claims = jwt.MapClaims{}
	}

	if cfg.CtxKey == "" {
		cfg.CtxKey = "user"
	}

	if cfg.TFR == "" {
		cfg.TFR = "header:" + fiber.HeaderAuthorization
	}

	if cfg.AuthScheme == "" {
		cfg.AuthScheme = "Bearer"
	}

	cfg.keyFunc = func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != cfg.Method {
			return nil, fmt.Errorf("unexpected signing method %v", t.Header["alg"])
		}
		if len(cfg.Keys) > 0 {
			if kid, ok := t.Header["kid"].(string); ok {
				if key, ok := cfg.Keys[kid]; ok {
					return key, nil
				}
			}
			return nil, fmt.Errorf("unexpected jwt key id %v", t.Header["kid"])

		}
		return cfg.Key, nil
	}

	parts := strings.Split(cfg.TFR, ":")
	header := Header(parts[1], cfg.AuthScheme)
	switch parts[0] {
	case "query":
		header = Query(parts[1])
	case "param":
		header = Param(parts[1])
	case "cookie":
		header = Cookie(parts[1])
	}

	return func(c *fiber.Ctx) error {
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()

		}

		auth, err := header(c)
		if err != nil {
			cfg.Error(c, err)
			return nil
		}

		token := new(jwt.Token)
		if _, ok := cfg.Claims.(jwt.MapClaims); ok {
			token, err = jwt.Parse(auth, cfg.keyFunc)
		} else {
			t := reflect.ValueOf(cfg.Claims).Type().Elem()
			claims := reflect.New(t).Interface().(jwt.Claims)
			token, err = jwt.ParseWithClaims(auth, claims, cfg.keyFunc)
		}

		if err == nil && token.Valid {
			c.Locals(cfg.CtxKey, token)
			cfg.Successful(c)
			return c.Next()
		}
		cfg.Error(c, err)
		return nil
	}

}

func Cookie(name string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Cookies(name)
		if token == "" {
			return "", fmt.Errorf("missing or malformed jwt")
		}
		return token, nil
	}
}

func Header(header, authScheme string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		auth := c.Get(header)
		length := len(authScheme)
		if len(auth) > length+1 && auth[:length] == authScheme {
			return auth[length+1:], nil
		}

		return "", fmt.Errorf("missing or malformed jwt")
	}
}

func Param(param string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Params(param)
		if token == "" {
			return "", fmt.Errorf("missing or malformed jwt")
		}
		return token, nil
	}
}

func Query(query string) func(c *fiber.Ctx) (string, error) {
	return func(c *fiber.Ctx) (string, error) {
		token := c.Query(query)
		if token == "" {
			return "", fmt.Errorf("missing or malformed jwt")
		}
		return token, nil
	}
}
