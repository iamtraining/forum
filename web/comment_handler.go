package web

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
)

type CommentHandler struct {
	store *store.Store
}

func (h *CommentHandler) createComment(c *fiber.Ctx) error {
	idstr := c.Params("postID")
	id, err := uuid.Parse(idstr)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse postID",
		})
		return nil
	}

	form := CreateCommentForm{
		Content: `json:"content"`,
	}

	err = c.BodyParser(&form)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return nil
	}

	if !form.Validate() {
		return c.Redirect(c.Request().URI().String(), fiber.StatusFound)
	}

	post, err := h.store.ReadPost(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting post",
		})
		return nil
	}

	if err := h.store.CreateComment(&entity.ForumComment{
		ID:      uuid.New(),
		PostID:  post.ID,
		Content: form.Content,
	}); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": " failure while creating comment " + err.Error(),
		})
	} else {
		return c.Redirect("/threads/"+post.ThreadID.String()+"/"+idstr, fiber.StatusFound)
	}
	return nil
}

func (h *CommentHandler) deleteComment(c *fiber.Ctx) error {
	type data struct {
		ID string `json:"comment_id"`
	}

	var body data

	err := c.BodyParser(&body)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "cannot parse request body",
		})
		return nil
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse comment id",
		})
		return nil
	}

	if err = h.store.DeleteComment(id); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while deleting post " + err.Error(),
		})
	} else {
		c.Send([]byte("comment deleted"))
	}
	return nil
}

func (h *CommentHandler) updateComment(c *fiber.Ctx) error {
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
		return nil
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	comment, err := h.store.ReadComment(id)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting comment " + err.Error(),
		})
		return nil
	}

	if body.Content == "" {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "you didnt fill the required fields(content)",
		})
		return nil
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
			return nil
		}
		c.Send(j)
	}

	return nil
}
