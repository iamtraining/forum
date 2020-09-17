package web

import (
	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
)

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

func (h *Handler) Home() func(*fiber.Ctx) {
	type data struct {
		SessionData
		Threads []entity.ForumThread
	}
	return func(c *fiber.Ctx) {
		threads, err := h.store.Threads()
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while getting threads",
			})
			return
		}

		if err := c.Render("home", data{
			SessionData: GetSessionData(c),
			Threads:     threads,
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) Post() func(*fiber.Ctx) {
	type data struct {
		SessionData
		Thread   entity.ForumThread
		Post     entity.ForumPost
		Comments []entity.ForumComment
	}
	return func(c *fiber.Ctx) {
		tID, err := uuid.Parse(c.Params("threadID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return
		}

		pID, err := uuid.Parse(c.Params("postID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse post id",
			})
			return
		}

		thread, err := h.store.ReadThread(tID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return
		}

		post, err := h.store.ReadPost(pID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting post",
			})
			return
		}

		comments, err := h.store.ReadCommentsByPost(post.ID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting comments",
			})
			return
		}

		if err := c.Render("post", data{
			SessionData: GetSessionData(c),
			Thread:      thread,
			Post:        post,
			Comments:    comments,
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) CreatePost() func(*fiber.Ctx) {
	type data struct {
		SessionData
		Thread entity.ForumThread
	}
	return func(c *fiber.Ctx) {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return
		}

		thread, err := h.store.ReadThread(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return
		}

		if err := c.Render("post_create", data{
			SessionData: GetSessionData(c),
			Thread:      thread,
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) Thread() func(*fiber.Ctx) {
	type data struct {
		SessionData
		Thread entity.ForumThread
		Posts  []entity.ForumPost
	}
	return func(c *fiber.Ctx) {
		id, err := uuid.Parse(c.Params("threadID"))
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "cannot parse thread id",
			})
			return
		}

		thread, err := h.store.ReadThread(id)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return
		}

		posts, err := h.store.ReadPostsByThread(thread.ID)
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting posts",
			})
			return
		}

		if err := c.Render("thread", data{
			SessionData: GetSessionData(c),
			Thread:      thread,
			Posts:       posts,
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) CreateThread() func(*fiber.Ctx) {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) {
		if err := c.Render("thread_create", data{
			SessionData: GetSessionData(c),
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) ThreadList() func(*fiber.Ctx) {
	type data struct {
		SessionData
		Threads []entity.ForumThread
	}
	return func(c *fiber.Ctx) {
		threads, err := h.store.Threads()
		if err != nil {
			c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "failure while getting thread",
			})
			return
		}

		if err := c.Render("threads", data{
			SessionData: GetSessionData(c),
			Threads:     threads,
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) LoginPage() func(*fiber.Ctx) {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) {
		if err := c.Render("login", data{
			SessionData: GetSessionData(c),
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}

func (h *Handler) RegisterPage() func(*fiber.Ctx) {
	type data struct {
		SessionData
	}
	return func(c *fiber.Ctx) {
		if err := c.Render("register", data{
			SessionData: GetSessionData(c),
		}); err != nil {
			c.Status(500).Send(err.Error())
			return
		}
	}
}
