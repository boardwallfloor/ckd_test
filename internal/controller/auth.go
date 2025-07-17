package controller

import (
	"boardwallfloor/ckd/internal/db"
	"boardwallfloor/ckd/internal/service"
	"boardwallfloor/ckd/internal/util"
	"encoding/json"
	"net/http"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(s service.AuthService) *AuthController {
	return &AuthController{authService: s}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	params := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := c.authService.Register(r.Context(), params)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, util.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	token, err := c.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		util.RespondWithError(w, http.StatusUnauthorized, "Login failed: "+err.Error())
		return
	}

	util.RespondWithJSON(w, http.StatusOK, util.Response{
		Success: true,
		Message: "Login successful",
		Data:    map[string]string{"token": token},
	})
}
