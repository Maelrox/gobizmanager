// internal/auth/handler.go
package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"gobizmanager/internal/app/user"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/logger"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	UserRepo   *user.Repository
	JWTManager *JWTManager
	Validator  *validator.Validate
	MsgStore   *language.MessageStore
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req user.CreateUserRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		utils.RespondError(w, r, h.MsgStore, err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		logger.Error("Validation failed", zap.Error(err))
		utils.ValidationError(w, r, err, h.MsgStore)
		return
	}

	// Check if username already exists
	_, err := h.UserRepo.GetUserByEmail(req.Username)
	if err == nil {
		logger.Info("Username already exists", zap.String("username", req.Username))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthUsernameExists))
		return
	}

	// Register user
	userID, err := h.UserRepo.RegisterUser(req.Username, req.Password, req.Phone)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthCreateUserFailed))
		return
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(userID)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthTokenGenFailed))
		return
	}

	logger.Info("New user registered successfully", zap.Int64("userID", userID))
	utils.JSON(w, http.StatusCreated, tokens)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRequest))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthValidationFailed))
		return
	}

	// Get user by username
	u, err := h.UserRepo.GetUserByEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidCredentials))
		} else {
			utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthDatabaseError))
		}
		return
	}

	// Check password
	if !user.CheckPassword(req.Password, u.Password) {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidCredentials))
		return
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(u.ID)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req user.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRequest))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthValidationFailed))
		return
	}

	claims, err := h.JWTManager.VerifyToken(req.RefreshToken)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRefreshToken))
		return
	}

	_, err = h.UserRepo.GetUserByID(claims.UserID)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthUserNotFound))
		return
	}

	// Generate new token pair
	tokens, err := h.JWTManager.GenerateTokenPair(claims.UserID)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}
