package tryout

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/middleware"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetCurrentTryout(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	data, err := h.service.GetCurrentTryout(c.UserContext(), userID)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "current tryout fetched", data)
}

func (h *Handler) StartCurrentTryout(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	data, err := h.service.StartCurrentTryout(c.UserContext(), userID)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "attempt prepared", data)
}

func (h *Handler) GetCurrentAttempt(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	data, err := h.service.GetCurrentAttempt(c.UserContext(), userID)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "current attempt fetched", data)
}

func (h *Handler) SaveCurrentAnswers(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	var req SaveAnswersRequest
	if err := c.BodyParser(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}
	data, err := h.service.SaveCurrentAnswers(c.UserContext(), userID, req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "answers saved", data)
}

func (h *Handler) SubmitCurrentAttempt(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	var req SubmitAttemptRequest
	if len(c.Body()) > 0 {
		if err := c.BodyParser(&req); err != nil {
			return httpx.BadRequest("invalid request body")
		}
	}
	data, err := h.service.SubmitCurrentAttempt(c.UserContext(), userID, req)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "attempt submitted", data)
}

func (h *Handler) GetCurrentResult(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	data, err := h.service.GetCurrentResult(c.UserContext(), userID)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "result fetched", data)
}

func (h *Handler) GetCurrentLeaderboard(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return httpx.Unauthorized("invalid token context")
	}

	userID := claims.UserID
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	data, err := h.service.GetCurrentLeaderboard(c.UserContext(), userID, limit, offset)
	if err != nil {
		return err
	}
	return httpx.Success(c, fiber.StatusOK, "leaderboard fetched", data)
}
