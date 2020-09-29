package web

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
)

type PostHandler struct {
	store *store.Store
}

func (h *PostHandler) getPost(c *fiber.Ctx) error {
	type data struct {
		ThreadID string `json:"thread_id"`
		PostID   string `json:"post_id"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return nil
	}

	threadID, err := uuid.Parse(body.ThreadID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse thread id",
		})
		return nil
	}

	postID, err := uuid.Parse(body.PostID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse post id",
		})
		return nil
	}

	thread, err := h.store.ReadThread(threadID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return nil
	}

	post, err := h.store.ReadPost(postID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting post",
		})
		return nil
	}

	comments, err := h.store.ReadCommentsByPost(post.ID)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting comments " + err.Error(),
		})
		return nil
	}

	c.Status(fiber.StatusOK).JSON(fiber.Map{
		"thread_id":   thread.ID,
		"title":       thread.Title,
		"description": thread.Description,
		"posts": struct {
			ID       uuid.UUID `json"post_id"`
			Title    string    `json"title"`
			Content  string    `json"content"`
			Count    int       `json"count"`
			Comments []entity.ForumComment
		}{
			ID:       post.ID,
			Title:    post.Title,
			Content:  post.Content,
			Count:    post.Count,
			Comments: comments,
		},
	})

	return nil
}

func (h *PostHandler) getPostsByThread(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse thread id",
		})
		return nil
	}
	/*type data struct {
		ID string `json:"id"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return
	}


	id, err := uuid.Parse(body.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse thread id",
		})
		return
	}
	*/

	posts, err := h.store.ReadPostsByThread(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "faulure while getting threads",
		})
		return nil
	}

	j, err := json.Marshal(&posts)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while marshalling threads",
		})
		return nil
	}

	return c.Send(j)
}

func (h *PostHandler) createPost(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	form := CreatePostForm{
		Title:   `json:"title"`,
		Content: `json:"content"`,
	}
	err = c.BodyParser(&form)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return nil
	}

	thread, err := h.store.ReadThread(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return nil
	}

	post := &entity.ForumPost{
		ID:       uuid.New(),
		ThreadID: thread.ID,
		Title:    form.Title,
		Content:  form.Content,
	}

	err = h.store.CreatePost(post)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": " failure while creating post",
		})
		return nil
	}

	return c.Redirect("/threads/"+thread.ID.String()+"/"+post.ID.String(), fiber.StatusFound)
}

func (h *PostHandler) deletePost(c *fiber.Ctx) error {
	type data struct {
		PostID string `json:"id"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return nil
	}

	id, err := uuid.Parse(body.PostID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse post id",
		})
		return nil
	}

	if err = h.store.DeletePost(id); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while deleting post",
		})
	} else {
		c.Send([]byte("post deleted"))
	}

	return nil

}

func (h *PostHandler) updatePost(c *fiber.Ctx) error {
	type data struct {
		ID      string `json:"post_id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	post, err := h.store.ReadPost(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting post",
		})
		return nil
	}

	if body.Title == "" && body.Content == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "you didnt fill the required fields(title and content)",
		})
		return nil
	}

	if body.Title == "" && body.Content != "" {
		post.Content = body.Content
	}

	if body.Title != "" && body.Content != "" {
		post.Title, post.Content = body.Title, body.Content
	}

	if err := h.store.UpdatePost(&post); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while updating post",
		})
	} else {
		j, err := json.Marshal(&post)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while marshalling post",
			})
			return nil
		}
		c.Send(j)
	}

	return nil
}
