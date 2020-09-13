package web

import (
	"github.com/gofiber/fiber"
	"github.com/gofiber/fiber/middleware"
	"github.com/gofiber/session"
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
		App:   fiber.New(),
	}

	threads := ThreadHandler{store: store}
	posts := PostHandler{store: store}
	comments := CommentHandler{store: store}
	users := UserHandler{store: store}

	h.App.Use(middleware.Logger())
	h.App.Use(IsAuthed)

	h.App.Post("/register", users.Register)
	h.App.Post("/login", users.Login)
	h.App.Post("logout", users.Logout)

	routes := h.App.Group("/threads")
	routes.Post("/", threads.getThread)
	routes.Get("/", threads.getThreads)
	routes.Patch("/", threads.updateThread)
	routes.Delete("/delete", threads.deleteThread)
	routes.Post("/new", threads.createThread)

	routes.Post("/posts", posts.getPost)
	routes.Get("/posts/:id", posts.getPostsByThread)
	routes.Patch("/posts", posts.updatePost)
	routes.Delete("/posts/delete", posts.deletePost)
	routes.Post("/posts/new", posts.createPost)

	routes.Post("/posts/comments", comments.createComment)
	routes.Delete("/posts/comments/delete", comments.deleteComment)
	routes.Patch("/posts/comments", comments.updateComment)

	return h
}
