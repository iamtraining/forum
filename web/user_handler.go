package web

import (
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	store *store.Store
}

func (h *UserHandler) Register(c *fiber.Ctx) {
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
		return
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
		st := s.Get(c)
		st.Set("form", form)
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return
	}

	if err := h.store.Create(&entity.User{
		ID:       uuid.New(),
		Username: form.Username,
		Password: string(password),
	}); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "your registration was successful. please log in",
		"username": form.Username,
		"password": form.Password,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) {
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
		return
	}

	form := LoginForm{
		Username:             body.Username,
		Password:             body.Password,
		IncorrectCredentials: false,
	}

	if !form.Validate() {
		st := s.Get(c)
		st.Set("form", form)
		return
	}

	user, err := h.store.GetUserByUsername(form.Username)
	if err != nil {
		form.IncorrectCredentials = true

		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": form.Err,
		})
	} else {
		pwErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password))
		form.IncorrectCredentials = pwErr != nil
	}

	st := s.Get(c)
	check := st.Get("islogined")
	if check == nil {
		st.Set("islogined", true)
		st.Save()
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "you have been logged in successfully",
		})
	} else {
		c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "you are already logged in",
		})
	}

}

func (h *UserHandler) Logout(c *fiber.Ctx) {
	store := s.Get(c)
	store.Delete("islogined")
	err := store.Save()
	if err != nil {
		panic(err)
	}
	c.ClearCookie()
	c.Send("you are logged out")
	return
}
