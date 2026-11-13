package handler

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"redditclone/internal/user"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/golang-jwt/jwt/v5"
)

type AuthInterface interface {
	Create(string, string) (user.User, error)
	Login(string, string) (user.User, error)
}

type AuthHandler struct {
	users  AuthInterface
	secret []byte
}

func NewAuthHandler(users AuthInterface, secret string) *AuthHandler {
	return &AuthHandler{users: users, secret: []byte(secret)}
}

type LoginForm struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	User user.User `json:"user"`
	jwt.RegisteredClaims
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	var form *LoginForm
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := decoder.Decode(&form); err != nil || form == nil {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	validUsername, _ := regexp.MatchString(`^[A-Za-z0-9_-]{1,32}$`, form.Username)
	if !validUsername {
		errorValidation(w, 422, "username", "bad username")
		return
	}
	if utf8.RuneCountInString(form.Password) < 8 || len(form.Password) > 72 {
		errorValidation(w, 422, "password", "bad password")
		return
	}
	u, err := h.users.Create(form.Username, form.Password)
	if errors.Is(err, user.ErrorUsername) {
		errorValidation(w, 422, "username", "already exists")
		return
	}
	if err != nil {
		log.Printf("error create user %v", err)
		errorJSON(w, 500, map[string]string{"message": "internal server error"})
		return
	}

	now := time.Now()
	claims := &Claims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.secret)
	if err != nil {
		log.Printf("error sign user %v", err)
		errorJSON(w, 500, map[string]string{"message": "internal server error"})
		return
	}
	writeJSON(w, 201, map[string]string{"token": tokenString})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	log.Printf("Method: %s, Path: %s \n", r.Method, r.URL.Path)
	var form *LoginForm
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	if err := decoder.Decode(&form); err != nil || form == nil {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		errorJSON(w, 400, "cannot decode payload")
		return
	}
	validateUsername, _ := regexp.MatchString(`^[A-Za-z_-]{1,32}$`, form.Username)
	if !validateUsername {
		errorValidation(w, 422, "username", "bad username")
		return
	}
	if utf8.RuneCountInString(form.Password) < 8 || len(form.Password) > 72 {
		errorValidation(w, 422, "password", "bad password")
		return
	}
	u, err := h.users.Login(form.Username, form.Password)
	if errors.Is(err, user.ErrorUnknownUser) {
		errorValidation(w, 401, "login", "unknown username or password")
		return
	}
	now := time.Now()
	claims := &Claims{
		User: u,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(7 * 24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(h.secret)
	if err != nil {
		log.Printf("error sign user %v", err)
		errorJSON(w, 500, map[string]string{"message": "internal server error"})
		return
	}
	writeJSON(w, 200, map[string]string{"token": tokenString})
}

func errorJSON(w http.ResponseWriter, status int, mes any) {
	writeJSON(w, status, map[string]any{
		"status": status,
		"error":  mes,
	})
}

func errorValidation(w http.ResponseWriter, status int, param string, msg string) {
	writeJSON(w, status, map[string]interface{}{
		"errors": []map[string]string{
			{
				"location": "body",
				"param":    param,
				"msg":      msg,
			},
		},
	})
}

func writeJSON(w http.ResponseWriter, status int, mes any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(mes); err != nil {
		log.Printf("encode response: %v", err)
	}
}
