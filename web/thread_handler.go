package web

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
)

type ThreadHandler struct {
	store *store.Store
}

func (h *ThreadHandler) getThread(ctx *fiber.Ctx) error {
	type data struct {
		ID string `json:"id"`
	}

	var body data
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	thread, err := h.store.ReadThread(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return nil
	}

	j, err := json.Marshal(&thread)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while marshalling thread",
		})
		return nil
	}

	return ctx.Send(j)
}

func (h *ThreadHandler) getThreads(ctx *fiber.Ctx) error {
	threads, err := h.store.Threads()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "faulure while getting threads",
		})
		return nil
	}

	j, err := json.Marshal(&threads)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while marshalling threads",
		})
		return nil
	}

	return ctx.Send(j)
}

func (h *ThreadHandler) createThread(ctx *fiber.Ctx) error {
	form := CreateThreadForm{
		Title:       `json:"title"`,
		Description: `json:"description"`,
	}
	err := ctx.BodyParser(&form)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	if err := h.store.CreateThread(&entity.ForumThread{
		ID:          uuid.New(),
		Title:       form.Title,
		Description: form.Description,
	}); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while creating thread",
		})
	}

	return ctx.Redirect("/threads", fiber.StatusFound)
}

func (h *ThreadHandler) updateThread(ctx *fiber.Ctx) error {
	type data struct {
		ID          string `json:"id"`
		Title       string `json:"title,omitempty"`
		Description string `json:"description,omitempty"`
	}

	var body data

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return nil
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	curr, err := h.store.ReadThread(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return nil
	}

	if body.Title == "" && body.Description == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "you didnt fill the required fields(title and description)",
		})
		return nil
	}

	if body.Title == "" && body.Description != "" {
		curr.Description = body.Description
	}

	if body.Title != "" && body.Description != "" {
		curr.Title, curr.Description = body.Title, body.Description
	}

	if err := h.store.UpdateThread(&curr); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while updating thread",
		})
	} else {
		j, err := json.Marshal(&curr)
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failure while marshalling threads",
			})
			return nil
		}
		ctx.Send(j)
	}

	return nil
}

func (h *ThreadHandler) deleteThread(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return nil
	}

	if err = h.store.DeleteThread(id); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while delete thread",
		})
	} else {
		ctx.Redirect("/threads", fiber.StatusFound)
	}

	return nil
}
