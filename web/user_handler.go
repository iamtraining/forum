package web

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/session/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
	"golang.org/x/crypto/bcrypt"
)

var sessions *session.Session

type UserHandler struct {
	store *store.Store
}

func init() {
	sessions = session.New()
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	type data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	body := data{}
	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	form := CreateUserForm{
		Username:       body.Username,
		Password:       body.Password,
		IsNotAvailable: false,
	}

	if _, err := h.store.GetUserByUsername(form.Username); err == nil {
		form.IsNotAvailable = true
	}

	if !form.Validate() {
		c.Locals("form", form)
		return c.Redirect("/register", fiber.StatusFound)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	if err := h.store.Create(&entity.User{
		ID:       uuid.New(),
		Username: form.Username,
		Password: string(password),
	}); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return nil
	}

	return c.Redirect("/", fiber.StatusFound)

	/*c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "your registration was successful. please log in",
		"username": form.Username,
		"password": form.Password,
	})
	*/

}

func (h *UserHandler) PrepareLogin(c *fiber.Ctx) error {
	type data struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	body := data{}
	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	form := LoginForm{
		Username:             body.Username,
		Password:             body.Password,
		IncorrectCredentials: false,
	}

	user, err := h.store.GetUserByUsername(form.Username)
	if err != nil {
		form.IncorrectCredentials = true
		fmt.Println("invalid credentials")
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": form.Err,
		})
	} else {
		pwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		form.IncorrectCredentials = pwErr != nil
		fmt.Println("valid credentials")
	}

	if !form.Validate() {
		c.Locals("form", form)
		return c.Redirect("/login", fiber.StatusFound)
	}
	fmt.Println(form)

	c.Locals("user", user)

	return c.Next()
}

func (h *UserHandler) CommitLogin(c *fiber.Ctx) error {
	user := c.Locals("user").(entity.User)

	t, err := Login(c, user.ID, []byte("SECRET_KEY"))
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "commit login error",
		})
	}

	fmt.Println("login t", t)

	return c.Redirect("/", fiber.StatusFound)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	Logout(c)

	return c.Redirect("/", fiber.StatusFound)
}
