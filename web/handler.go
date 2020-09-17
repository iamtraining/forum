package web

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/session"
	"github.com/gofiber/template/html"
	"github.com/iamtraining/forum/store"
)

type Handler struct {
	store   store.Store
	App     *fiber.App
	Session *session.Session
}

func NewHandler(store *store.Store) *Handler {
	h := &Handler{
		store: *store,
		App: fiber.New(&fiber.Settings{
			Views: html.New("./templates/", ".html"),
		}),
	}

	threads := ThreadHandler{store: store}
	posts := PostHandler{store: store}
	comments := CommentHandler{store: store}
	users := UserHandler{store: store}

	h.App.Use(middleware.Logger())
	h.App.Use(h.IsLoggedIn)
	//h.App.Use(csrf.New())

	/*h.App.Get("/", func(c *fiber.Ctx) {
		c.Send(c.Locals("csrf"))
	})
	*/

	h.App.Get("/", h.Home())
	h.App.Post("/register", users.Register)
	h.App.Post("/login", users.Login)
	h.App.Post("/logout", users.Logout)
	h.App.Get("/login", h.LoginPage())
	h.App.Get("/register", h.RegisterPage())
	h.App.Post("/thread/:id/delete", threads.deleteThread)

	routes := h.App.Group("/threads")
	//routes.Post("/", threads.getThread)
	//routes.Get("/", threads.getThreads)
	routes.Patch("/", threads.updateThread)
	routes.Post("/:id/delete", threads.deleteThread)
	routes.Post("/", threads.createThread)

	routes.Post("/posts", posts.getPost)
	routes.Get("/posts/:id", posts.getPostsByThread)
	routes.Patch("/posts", posts.updatePost)
	routes.Delete("/posts/delete", posts.deletePost)
	routes.Post("/:id", posts.createPost)

	routes.Post("/:threadID/:postID", comments.createComment)
	routes.Delete("/posts/comments/delete", comments.deleteComment)
	routes.Patch("/posts/comments", comments.updateComment)

	routes.Get("/", h.ThreadList())
	routes.Get("/:id/new", h.CreatePost())
	routes.Get("/:threadID/:postID", h.Post())
	routes.Get("/new", h.CreateThread())
	routes.Get("/:threadID", h.Thread())

	return h
}
