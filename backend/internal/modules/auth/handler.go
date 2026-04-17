package auth

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/httpx"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/middleware"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(router fiber.Router, authMW fiber.Handler) {
	authGroup := router.Group("/auth")

	authGroup.Post("/register", h.Register)
	authGroup.Post("/login", h.Login)
	authGroup.Post("/refresh", h.RefreshToken)

	protectedAuth := authGroup.Group("/", authMW)
	protectedAuth.Post("/logout", h.Logout)
	protectedAuth.Post("/logout-all", h.LogoutAll)

	router.Get("/me", authMW, h.Me)
}

// Register godoc
// @Summary Register user baru
// @Description Membuat akun user baru menggunakan email, nama lengkap, dan password.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body RegisterRequest true "Register payload"
// @Success 201 {object} UserSuccessResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 409 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c *fiber.Ctx) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid request body"})
	}

	resp, err := h.svc.Register(c.UserContext(), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrUserExists):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"success": false, "message": "email is already registered"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	return httpx.Success(c, fiber.StatusCreated, "user registered", resp)
}

// Login godoc
// @Summary Login user
// @Description Melakukan autentikasi user, membuat session, lalu menulis access token dan refresh token ke cookie HTTP-only.
// @Tags Auth
// @Accept json
// @Produce json
// @Param payload body LoginRequest true "Login payload"
// @Success 200 {object} LoginSuccessResponse
// @Failure 400 {object} httpx.ErrorResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"success": false, "message": "invalid request body"})
	}

	meta := SessionMeta{
		IPAddress: realIP(c),
		UserAgent: c.Get("User-Agent"),
		DeviceID:  c.Get("X-Device-ID"),
	}

	resp, accToken, refToken, err := h.svc.Login(c.UserContext(), &req, meta)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCredentials):
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "email or password is incorrect"})
		case errors.Is(err, ErrAccountInactive):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "account is inactive"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	h.setAuthCookies(c, accToken, refToken, resp.AccessExpiresAt, resp.RefreshExpiresAt)

	return httpx.Success(c, fiber.StatusOK, "login successful", resp)
}

// RefreshToken godoc
// @Summary Refresh session login
// @Description Membuat access token dan refresh token baru menggunakan cookie `refresh_token` yang masih valid.
// @Tags Auth
// @Produce json
// @Success 200 {object} LoginSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 403 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *fiber.Ctx) error {

	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "missing refresh token"})
	}

	resp, accToken, newRefToken, err := h.svc.RefreshToken(c.UserContext(), refreshToken)
	if err != nil {
		switch {
		case errors.Is(err, ErrSessionRevoked), errors.Is(err, ErrSessionExpired):
			h.clearAuthCookies(c)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "session has expired or been revoked"})
		case errors.Is(err, ErrAccountInactive):
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"success": false, "message": "account is inactive"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
		}
	}

	h.setAuthCookies(c, accToken, newRefToken, resp.AccessExpiresAt, resp.RefreshExpiresAt)

	return httpx.Success(c, fiber.StatusOK, "session refreshed", resp)
}

// Logout godoc
// @Summary Logout dari device saat ini
// @Description Mencabut session aktif berdasarkan `session_id` dari access token lalu menghapus cookie auth.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpx.EmptySuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *Handler) Logout(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "missing auth context"})
	}

	if err := h.svc.Logout(c.UserContext(), claims.SessionID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	h.clearAuthCookies(c)
	return httpx.Success(c, fiber.StatusOK, "logout successful", nil)
}

// LogoutAll godoc
// @Summary Logout semua device
// @Description Mencabut seluruh session milik user yang sedang login lalu menghapus cookie auth saat ini.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} httpx.EmptySuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/auth/logout-all [post]
func (h *Handler) LogoutAll(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "missing auth context"})
	}

	if err := h.svc.LogoutAll(c.UserContext(), claims.UserID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	h.clearAuthCookies(c)
	return httpx.Success(c, fiber.StatusOK, "logout successful", nil)
}

// Me godoc
// @Summary Ambil profil user login
// @Description Mengembalikan profil user berdasarkan subject di access token.
// @Tags Auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserSuccessResponse
// @Failure 401 {object} httpx.ErrorResponse
// @Failure 404 {object} httpx.ErrorResponse
// @Failure 500 {object} httpx.ErrorResponse
// @Router /api/v1/me [get]
func (h *Handler) Me(c *fiber.Ctx) error {
	claims, ok := middleware.ClaimsFromContext(c)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"success": false, "message": "missing auth context"})
	}

	resp, err := h.svc.Me(c.UserContext(), claims.UserID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"success": false, "message": "user not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"success": false, "message": err.Error()})
	}

	return httpx.Success(c, fiber.StatusOK, "profile fetched", resp)
}

func (h *Handler) setAuthCookies(c *fiber.Ctx, access string, refresh string, accExp time.Time, refExp time.Time) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    access,
		Expires:  accExp,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Expires:  refExp,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
}

func (h *Handler) clearAuthCookies(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
	})
}

func realIP(c *fiber.Ctx) string {
	if ip := c.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := c.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return c.IP()
}
