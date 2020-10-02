package web

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/session/v2"
	"github.com/gofiber/template/html"
	"github.com/iamtraining/forum/apiserver"
	"github.com/iamtraining/forum/store"
)

type Handler struct {
	store   store.Store
	App     *fiber.App
	Session *session.Session
}

func NewHandler(store *store.Store, cfg apiserver.Config) *Handler {
	engine := html.New("./templates/", ".html")

	h := &Handler{
		store: *store,
		App: fiber.New(fiber.Config{
			Views: engine,
		}),
	}

	threads := ThreadHandler{store: store}
	posts := PostHandler{store: store}
	comments := CommentHandler{store: store}
	users := UserHandler{store: store}

	h.App.Use(logger.New(logger.Config{
		Next:       nil,
		Format:     "[${time}] ${status} - ${latency} - ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
		Output:     os.Stderr,
	}))

	web := h.App.Group("")

	web.Get("/", h.Home())
	web.Post("/login", users.PrepareLogin, users.CommitLogin)
	web.Get("/login", h.LoginPage())
	web.Post("/register", users.Register)
	web.Get("/logout", users.Logout)
	web.Get("/register", h.RegisterPage())

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
