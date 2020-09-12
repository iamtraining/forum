package web

import (
	"encoding/json"

	"github.com/gofiber/fiber"
	"github.com/google/uuid"
	"github.com/iamtraining/forum/entity"
	"github.com/iamtraining/forum/store"
)

type ThreadHandler struct {
	store *store.Store
}

func (h *ThreadHandler) getThread(ctx *fiber.Ctx) {
	type data struct {
		ID string `json:"id"`
	}

	var body data
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	thread, err := h.store.ReadThread(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return
	}

	j, err := json.Marshal(&thread)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while marshalling thread",
		})
		return
	}

	ctx.Send(j)
}

func (h *ThreadHandler) getThreads(ctx *fiber.Ctx) {
	threads, err := h.store.Threads()
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "faulure while getting threads",
		})
		return
	}

	j, err := json.Marshal(&threads)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while marshalling threads",
		})
		return
	}

	ctx.Send(j)
}

func (h *ThreadHandler) createThread(ctx *fiber.Ctx) {
	type data struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	var body data
	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return
	}

	if err := h.store.CreateThread(&entity.ForumThread{
		ID:          uuid.New(),
		Title:       body.Title,
		Description: body.Description,
	}); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while creating thread",
		})
	} else {
		ctx.Send("thread created")
	}
}

func (h *ThreadHandler) updateThread(ctx *fiber.Ctx) {
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
		return
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	curr, err := h.store.ReadThread(id)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while getting thread",
		})
		return
	}

	if body.Title == "" && body.Description == "" {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "you didnt fill the required fields(title and description)",
		})
		return
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
			return
		}
		ctx.Send("thread updated " + string(j))
	}
}

func (h *ThreadHandler) deleteThread(ctx *fiber.Ctx) {
	type data struct {
		ID string `json"id"`
	}

	var body data

	err := ctx.BodyParser(&body)
	if err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while parsing params to a struct",
		})
		return
	}

	id, err := uuid.Parse(body.ID)
	if err != nil {
		ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse id",
		})
		return
	}

	if err = h.store.DeleteThread(id); err != nil {
		ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failure while delete thread",
		})
	} else {
		ctx.Send("thread deleted")
	}
}
