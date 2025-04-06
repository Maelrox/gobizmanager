// internal/auth/handler.go
package auth

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"gobizmanager/internal/logger"
	"gobizmanager/internal/user"
	"gobizmanager/pkg/context"
	"gobizmanager/pkg/language"
	"gobizmanager/pkg/utils"
)

type Handler struct {
	UserRepo   *user.Repository
	JWTManager *JWTManager
	Validator  *validator.Validate
	MsgStore   *language.MessageStore
}

func NewHandler(userRepo *user.Repository, jwtManager *JWTManager, msgStore *language.MessageStore) *Handler {
	return &Handler{
		UserRepo:   userRepo,
		JWTManager: jwtManager,
		Validator:  validator.New(),
		MsgStore:   msgStore,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req user.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Failed to decode request", zap.Error(err))
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidRequest))
		return
	}

	// Validate the request
	if err := h.Validator.Struct(req); err != nil {
		// Log the error first
		logger.Error("Validation failed", zap.Error(err))
		// Then handle the validation error
		utils.ValidationError(w, err, lang, h.MsgStore)
		return
	}

	// Check if username already exists
	_, err := h.UserRepo.GetUserByEmail(req.Username)
	if err == nil {
		logger.Info("Username already exists", zap.String("username", req.Username))
		utils.JSONError(w, http.StatusConflict, h.MsgStore.GetMessage(lang, language.MsgAuthUsernameExists))
		return
	}

	// Register user
	userID, err := h.UserRepo.RegisterUser(req.Username, req.Password, req.Phone)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgAuthCreateUserFailed))
		return
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(userID)
	if err != nil {
		logger.Error("Failed to generate tokens", zap.Error(err))
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgAuthTokenGenFailed))
		return
	}

	logger.Info("New user registered successfully", zap.Int64("userID", userID))
	utils.JSON(w, http.StatusCreated, tokens)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req user.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidRequest))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgAuthValidationFailed))
		return
	}

	// Get user by username
	u, err := h.UserRepo.GetUserByEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidCredentials))
		} else {
			utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgAuthDatabaseError))
		}
		return
	}

	// Check password
	if !user.CheckPassword(req.Password, u.Password) {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidCredentials))
		return
	}

	// Generate tokens
	tokens, err := h.JWTManager.GenerateTokenPair(u.ID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgAuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	lang := context.GetLanguage(r.Context())

	var req user.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidRequest))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.JSONError(w, http.StatusBadRequest, h.MsgStore.GetMessage(lang, language.MsgAuthValidationFailed))
		return
	}

	claims, err := h.JWTManager.VerifyToken(req.RefreshToken)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthInvalidRefreshToken))
		return
	}

	_, err = h.UserRepo.GetUserByID(claims.UserID)
	if err != nil {
		utils.JSONError(w, http.StatusUnauthorized, h.MsgStore.GetMessage(lang, language.MsgAuthUserNotFound))
		return
	}

	// Generate new token pair
	tokens, err := h.JWTManager.GenerateTokenPair(claims.UserID)
	if err != nil {
		utils.JSONError(w, http.StatusInternalServerError, h.MsgStore.GetMessage(lang, language.MsgAuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}
