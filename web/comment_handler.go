package web

import (
	"encoding/json"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
)

type CommentHandler struct {
	store *store.Store
}

func (h *CommentHandler) createComment(c *fiber.Ctx) {
	type data struct {
		ID      string `json:"post_id"`
		Content string `json:"content"`
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

	post, err := h.store.ReadPost(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting post",
		})
		return
	}

	if err := h.store.CreateComment(&entity.ForumComment{
		ID:      uuid.New(),
		PostID:  post.ID,
		Content: body.Content,
	}); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": " failure while creating comment " + err.Error(),
		})
	} else {
		c.Send("comment created")
	}
}

func (h *CommentHandler) deleteComment(c *fiber.Ctx) {
	type data struct {
		ID string `json:"comment_id"`
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
			"error": "cannot parse comment id",
		})
		return
	}

	if err = h.store.DeleteComment(id); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while deleting post " + err.Error(),
		})
	} else {
		c.Send("comment deleted")
	}
}

func (h *CommentHandler) updateComment(c *fiber.Ctx) {
	type data struct {
		ID      string `json:"comment_id"`
		Content string `json:"content"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	comment, err := h.store.ReadComment(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting comment " + err.Error(),
		})
		return
	}

	if body.Content == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "you didnt fill the required fields(content)",
		})
		return
	}

	if body.Content != "" {
		comment.Content = body.Content
	}

	if err := h.store.UpdateComment(&comment); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while updating comment",
		})
	} else {
		j, err := json.Marshal(&comment)
		if err != nil {
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while marshalling comment",
			})
			return
		}
		c.Send("comment updated " + string(j))
	}
}
