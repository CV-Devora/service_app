package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	v1 "toko-emas/api/v1"
	"toko-emas/internal/conf"
	"toko-emas/internal/data"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo         *data.UserRepo
	jwtSecret     []byte
	accessTokenTT time.Duration
}

type authClaims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Type     string `json:"type"`
	jwt.RegisteredClaims
}

func NewAuthService(repo *data.UserRepo, cfg *conf.Config) *AuthService {
	secret := cfg.Auth.JWTSecret
	if secret == "" {
		secret = os.Getenv("JWT_SECRET")
	}
	ttl := 6 * time.Hour
	if cfg.Auth.AccessTokenTTL != "" {
		if parsed, err := time.ParseDuration(cfg.Auth.AccessTokenTTL); err == nil {
			ttl = parsed
		}
	}
	return &AuthService{
		repo:         repo,
		jwtSecret:     []byte(secret),
		accessTokenTT: ttl,
	}
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var req v1.AuthLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}

	user, err := s.repo.FindByUsername(req.Username)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "invalid credentials"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "invalid credentials"})
		return
	}

	resp, err := s.issueTokens(user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: resp})
}

func (s *AuthService) Refresh(w http.ResponseWriter, r *http.Request) {
	var req v1.AuthRefreshRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	tokenString := strings.TrimSpace(req.RefreshToken)
	if tokenString == "" {
		tokenString = bearerToken(r)
	}
	if tokenString == "" {
		writeJSON(w, http.StatusBadRequest, v1.Response{Code: 400, Message: "invalid body"})
		return
	}

	token, err := s.parseToken(tokenString)
	if err != nil || token.Claims == nil {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "invalid token"})
		return
	}

	claims, ok := token.Claims.(*authClaims)
	if !ok || claims.Type != "refresh" {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "invalid token"})
		return
	}

	user, err := s.repo.FindByUsername(claims.Username)
	if err != nil || user == nil {
		writeJSON(w, http.StatusUnauthorized, v1.Response{Code: 401, Message: "invalid token"})
		return
	}

	resp, err := s.issueTokens(user)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, v1.Response{Code: 500, Message: err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, v1.Response{Code: 200, Message: "success", Data: resp})
}

func (s *AuthService) issueTokens(user *data.User) (*v1.AuthTokenResponse, error) {
	accessToken, err := s.signToken(user, "access", s.accessTokenTT)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.signToken(user, "refresh", s.accessTokenTT)
	if err != nil {
		return nil, err
	}
	return &v1.AuthTokenResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		ExpiresIn: int64(s.accessTokenTT.Seconds()),
		TokenType: "Bearer",
		User: map[string]any{
			"id": user.ID,
			"nama": user.Nama,
			"username": user.Username,
			"role": user.Role,
		},
	}, nil
}

func (s *AuthService) signToken(user *data.User, tokenType string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := authClaims{
		UserID:   user.ID.String(),
		Username: user.Username,
		Role:     user.Role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tkn.SignedString(s.jwtSecret)
}

func (s *AuthService) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &authClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
}

func bearerToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
