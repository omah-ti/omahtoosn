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

// GetCurrentTryout godoc
// @Summary Ambil tryout aktif saat ini
// @Description Mengembalikan metadata tryout ongoing beserta ringkasan attempt user bila sudah ada.
// @Tags Tryout
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentTryoutSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/tryouts/current [get]
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

// StartCurrentTryout godoc
// @Summary Mulai attempt tryout aktif
// @Description Membuat atau mengambil attempt ongoing user untuk tryout aktif dan menyesuaikan expiry bila diperlukan.
// @Tags Tryout
// @Produce json
// @Security BearerAuth
// @Success 200 {object} StartCurrentTryoutSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/tryouts/current/start [post]
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

// GetCurrentAttempt godoc
// @Summary Ambil detail attempt aktif
// @Description Mengembalikan detail tryout, ringkasan attempt, daftar soal, jawaban tersimpan, dan navigator soal.
// @Tags Tryout
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentAttemptSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/attempts/current [get]
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

// SaveCurrentAnswers godoc
// @Summary Simpan jawaban attempt aktif
// @Description Menyimpan jawaban partial user untuk attempt aktif dan memperbarui statistik attempt.
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body SaveAnswersRequest true "Save answers payload"
// @Success 200 {object} SaveCurrentAnswersSuccessResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/attempts/current/answers [put]
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

// SubmitCurrentAttempt godoc
// @Summary Submit attempt aktif
// @Description Menyimpan final answers opsional, menghitung skor, lalu menutup attempt menjadi submitted atau auto_submitted.
// @Tags Tryout
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param payload body SubmitAttemptRequest false "Submit attempt payload"
// @Success 200 {object} SubmitCurrentAttemptSuccessResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/attempts/current/submit [post]
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

// GetCurrentResult godoc
// @Summary Ambil hasil tryout aktif
// @Description Mengembalikan ringkasan hasil current tryout untuk user login.
// @Tags Tryout
// @Produce json
// @Security BearerAuth
// @Success 200 {object} CurrentResultSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/results/current [get]
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

// GetCurrentLeaderboard godoc
// @Summary Ambil leaderboard tryout aktif
// @Description Mengembalikan leaderboard current tryout berikut data ranking user login bila tersedia.
// @Tags Tryout
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Jumlah data per halaman" default(50) minimum(1) maximum(100)
// @Param offset query int false "Offset pagination" default(0) minimum(0)
// @Success 200 {object} CurrentLeaderboardSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/leaderboard/current [get]
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
