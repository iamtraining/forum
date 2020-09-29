package web

import (
	"encoding/gob"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
)

func init() {
	gob.Register(uuid.UUID{})
}

/*func Render(c *fiber.Ctx, f1, f2 string, data interface{}) {
	tmpl := template.Must(template.ParseFiles(f1, f2))
	w := bytes.Buffer{}
	tmpl.ExecuteTemplate(&w, "", data)
	c.Render("layout", fiber.Map{
		"IsLoggedIn": false,
		"Data":       template.HTML(w.Bytes()),
	})
}

func (h *Handler) Home() func(*fiber.Ctx) {
	return func(c *fiber.Ctx) {
		threads, err := h.store.Threads()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while getting threads",
			})
			return
		}

		Render(c, "templates/new/layout.html", "templates/new/home.html", fiber.Map{
			"Threads": threads,
		})
	}
}
*/

func (h *Handler) Home() fiber.Handler {
	type data struct {
		SessionData
		Threads []entity.ForumThread
	}
	return func(c *fiber.Ctx) error {
		threads, err := h.store.Threads()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while getting threads",
			})
			return nil
		}

		return c.Render("home", data{
			SessionData: h.GetSessionData(c),
			Threads:     threads,
		})
	}
}

func (h *Handler) Post() fiber.Handler {
	type data struct {
		SessionData
		Thread   entity.ForumThread
		Post     entity.ForumPost
		Comments []entity.ForumComment
	}
	return func(c *fiber.Ctx) error {
		tID, err := uuid.Parse(c.Params("threadID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return nil
		}

		pID, err := uuid.Parse(c.Params("postID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse post id",
			})
			return nil
		}

		thread, err := h.store.ReadThread(tID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return nil
		}

		post, err := h.store.ReadPost(pID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting post",
			})
			return nil
		}

		comments, err := h.store.ReadCommentsByPost(post.ID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting comments",
			})
			return nil
		}

		return c.Render("post", data{
			SessionData: h.GetSessionData(c),
			Thread:      thread,
			Post:        post,
			Comments:    comments,
		})
	}
}

func (h *Handler) CreatePost() fiber.Handler {
	type data struct {
		SessionData
		Thread entity.ForumThread
	}
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return nil
		}

		thread, err := h.store.ReadThread(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return nil
		}

		return c.Render("post_create", data{
			SessionData: h.GetSessionData(c),
			Thread:      thread,
		})
	}
}

func (h *Handler) Thread() fiber.Handler {
	type data struct {
		SessionData
		Thread entity.ForumThread
		Posts  []entity.ForumPost
	}
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("threadID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return nil
		}

		thread, err := h.store.ReadThread(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return nil
		}

		posts, err := h.store.ReadPostsByThread(thread.ID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting posts",
			})
			return nil
		}

		return c.Render("thread", data{
			SessionData: h.GetSessionData(c),
			Thread:      thread,
			Posts:       posts,
		})
	}
}

func (h *Handler) CreateThread() fiber.Handler {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) error {
		return c.Render("thread_create", data{
			SessionData: h.GetSessionData(c),
		})
	}
}

func (h *Handler) ThreadList() fiber.Handler {
	type data struct {
		SessionData
		Threads []entity.ForumThread
	}
	return func(c *fiber.Ctx) error {
		threads, err := h.store.Threads()
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return nil
		}

		return c.Render("threads", data{
			SessionData: h.GetSessionData(c),
			Threads:     threads,
		})
	}
}

func (h *Handler) LoginPage() fiber.Handler {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) error {
		return c.Render("login", data{
			SessionData: h.GetSessionData(c),
		})
	}
}

func (h *Handler) RegisterPage() fiber.Handler {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) error {
		return c.Render("register", data{
			SessionData: h.GetSessionData(c),
		})
	}
}
