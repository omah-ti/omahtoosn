package auth

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/email"
	"github.com/omah-ti/omahtoosn/backend/internal/platform/security"
)

// error
var (
	ErrInvalidCredentials = errors.New("email or password is incorrect")
	ErrAccountInactive    = errors.New("account is inactive")
	ErrSessionExpired     = errors.New("session has expired")
	ErrSessionRevoked     = errors.New("session has been revoked")
	ErrPasswordTooWeak    = errors.New("password must be at least 8 characters")
)

// interface
type Service interface {
	Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error)
	Login(ctx context.Context, req *LoginRequest, meta SessionMeta) (*LoginResponse, string, string, error)
	ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error)
	ResetPassword(ctx context.Context, req *ResetPasswordRequest) error
	RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, string, string, error)
	Logout(ctx context.Context, sessionID string) error
	LogoutAll(ctx context.Context, userID string) error
	Me(ctx context.Context, userID string) (*UserResponse, error)
}

type SessionMeta struct {
	IPAddress string
	UserAgent string
	DeviceID  string
}

type PasswordResetConfig struct {
	FrontendURL string
	ResetPath   string
	TokenTTL    time.Duration
}

type service struct {
	repo              Repository
	tokens            security.TokenProvider
	emailSender       email.Sender
	passwordResetConf PasswordResetConfig
}

func NewService(repo Repository, tokens security.TokenProvider, emailSender email.Sender, passwordResetConf PasswordResetConfig) Service {
	if passwordResetConf.FrontendURL == "" {
		passwordResetConf.FrontendURL = "http://localhost:3000"
	}
	if passwordResetConf.ResetPath == "" {
		passwordResetConf.ResetPath = "/reset-password"
	}
	if passwordResetConf.TokenTTL <= 0 {
		passwordResetConf.TokenTTL = 30 * time.Minute
	}
	return &service{
		repo:              repo,
		tokens:            tokens,
		emailSender:       emailSender,
		passwordResetConf: passwordResetConf,
	}
}

// register
// flow: validasi → hash password → simpan user → return UserResponse
func (s *service) Register(ctx context.Context, req *RegisterRequest) (*UserResponse, error) {
	// Hash password sebelum disimpan
	hash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Email:        req.Email,
		FullName:     req.FullName,
		PasswordHash: hash,
		IsActive:     true,
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// login
// flow: cari user → verify password → cek aktif → buat session → return tokens
func (s *service) Login(ctx context.Context, req *LoginRequest, meta SessionMeta) (*LoginResponse, string, string, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, "", "", ErrInvalidCredentials
		}
		return nil, "", "", err
	}

	if !security.VerifyPassword(user.PasswordHash, req.Password) {
		return nil, "", "", ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, "", "", ErrAccountInactive
	}

	accToken, refToken, accExp, refExp, err := s.createSession(ctx, user, meta)
	if err != nil {
		return nil, "", "", err
	}

	_ = s.repo.UpdateLastLogin(ctx, user.ID)

	resp := &LoginResponse{
		AccessExpiresAt:  accExp,
		CookieStrategy:   "http_only",
		RefreshExpiresAt: refExp,
		User:             *toUserResponse(user),
	}
	return resp, accToken, refToken, nil
}

func (s *service) ForgotPassword(ctx context.Context, req *ForgotPasswordRequest) (*ForgotPasswordResponse, error) {
	emailAddress := strings.TrimSpace(req.Email)
	if emailAddress == "" {
		return nil, ErrUserNotFound
	}

	user, err := s.repo.GetUserByEmail(ctx, emailAddress)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if !user.IsActive {
		return nil, nil
	}

	resetToken, err := s.tokens.GenerateRefreshToken()
	if err != nil {
		return nil, err
	}

	expiresAt := time.Now().Add(s.passwordResetConf.TokenTTL)
	token := &PasswordResetToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: s.tokens.HashToken(resetToken),
		ExpiresAt: expiresAt,
	}
	if err := s.repo.CreatePasswordResetToken(ctx, token); err != nil {
		return nil, err
	}

	resetURL, err := s.buildPasswordResetURL(resetToken)
	if err != nil {
		return nil, err
	}
	if s.emailSender == nil {
		return nil, email.ErrNotConfigured
	}
	if err := s.emailSender.SendPasswordReset(ctx, email.PasswordResetMessage{
		To:        user.Email,
		FullName:  user.FullName,
		ResetURL:  resetURL,
		ExpiresAt: expiresAt,
	}); err != nil {
		return nil, err
	}

	return &ForgotPasswordResponse{}, nil
}

func (s *service) ResetPassword(ctx context.Context, req *ResetPasswordRequest) error {
	token := strings.TrimSpace(req.Token)
	if token == "" {
		return ErrPasswordResetTokenInvalid
	}
	if len(req.NewPassword) < 8 {
		return ErrPasswordTooWeak
	}

	hash, err := security.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.repo.ResetPasswordWithToken(ctx, s.tokens.HashToken(token), hash)
}

func (s *service) buildPasswordResetURL(token string) (string, error) {
	base, err := url.Parse(strings.TrimRight(s.passwordResetConf.FrontendURL, "/"))
	if err != nil {
		return "", err
	}
	path := s.passwordResetConf.ResetPath
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	rel, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	resetURL := base.ResolveReference(rel)
	query := resetURL.Query()
	query.Set("token", token)
	resetURL.RawQuery = query.Encode()
	return resetURL.String(), nil
}

func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, string, string, error) {
	tokenHash := s.tokens.HashToken(refreshToken)
	session, err := s.repo.GetSessionByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, ErrSessionNotFound) {
			return nil, "", "", ErrSessionRevoked
		}
		return nil, "", "", err
	}

	if session.RevokedAt != nil {
		return nil, "", "", ErrSessionRevoked
	}
	if time.Now().After(session.ExpiresAt) {
		return nil, "", "", ErrSessionExpired
	}

	user, err := s.repo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return nil, "", "", err
	}
	if !user.IsActive {
		return nil, "", "", ErrAccountInactive
	}

	if err := s.repo.RevokeSession(ctx, session.ID); err != nil {
		return nil, "", "", err
	}

	meta := SessionMeta{}
	if session.IPAddress != nil {
		meta.IPAddress = *session.IPAddress
	}
	if session.UserAgent != nil {
		meta.UserAgent = *session.UserAgent
	}
	if session.DeviceID != nil {
		meta.DeviceID = *session.DeviceID
	}

	accToken, refToken, accExp, refExp, err := s.createSession(ctx, user, meta)
	if err != nil {
		return nil, "", "", err
	}

	resp := &LoginResponse{
		AccessExpiresAt:  accExp,
		CookieStrategy:   "http_only",
		RefreshExpiresAt: refExp,
		User:             *toUserResponse(user),
	}
	return resp, accToken, refToken, nil
}

// cabut satu session berdasarkan sessionID dari JWT claims
func (s *service) Logout(ctx context.Context, sessionID string) error {
	sid, err := uuid.Parse(sessionID)
	if err != nil {
		return errors.New("invalid session id")
	}
	return s.repo.RevokeSession(ctx, sid)
}

// cabut semua session milik user (logout dari semua device)
func (s *service) LogoutAll(ctx context.Context, userID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return errors.New("invalid user id")
	}
	return s.repo.RevokeAllUserSessions(ctx, uid)
}

// ambil profil user yang sedang login
func (s *service) Me(ctx context.Context, userID string) (*UserResponse, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// generate tokens + simpan session ke DB
func (s *service) createSession(ctx context.Context, user *User, meta SessionMeta) (access string, refresh string, accExp time.Time, refExp time.Time, err error) {
	sessionID := uuid.New()

	access, err = s.tokens.GenerateAccessToken(user.ID.String(), user.Email, sessionID.String())
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	refresh, err = s.tokens.GenerateRefreshToken()
	if err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	accExp = time.Now().Add(s.tokens.AccessTokenTTL())
	refExp = time.Now().Add(s.tokens.RefreshTokenTTL())

	session := &AuthSession{
		ID:               sessionID,
		UserID:           user.ID,
		RefreshTokenHash: s.tokens.HashToken(refresh),
		ExpiresAt:        refExp,
	}

	if meta.IPAddress != "" {
		ip := meta.IPAddress
		session.IPAddress = &ip
	}
	if meta.UserAgent != "" {
		ua := meta.UserAgent
		session.UserAgent = &ua
	}
	if meta.DeviceID != "" {
		did := meta.DeviceID
		session.DeviceID = &did
	}

	if err = s.repo.CreateSession(ctx, session); err != nil {
		return "", "", time.Time{}, time.Time{}, err
	}

	return access, refresh, accExp, refExp, nil
}

func toUserResponse(u *User) *UserResponse {
	return &UserResponse{
		ID:          u.ID.String(),
		Email:       u.Email,
		FullName:    u.FullName,
		IsActive:    u.IsActive,
		LastLoginAt: u.LastLoginAt,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}
