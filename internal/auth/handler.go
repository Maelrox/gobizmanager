// internal/auth/handler.go
package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	types "gobizmanager/internal/types"
	user "gobizmanager/internal/user"
	"gobizmanager/pkg/encryption"
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

func NewHandler(userRepo *user.Repository, jwtManager *JWTManager, msgStore *language.MessageStore) *Handler {
	return &Handler{
		UserRepo:   userRepo,
		JWTManager: jwtManager,
		Validator:  validator.New(),
		MsgStore:   msgStore,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {

	var req types.CreateUserRequest
	if err := utils.ParseRequest(r, &req); err != nil {
		utils.RespondError(w, r, h.MsgStore, err)
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.ValidationError(w, r, err, h.MsgStore)
		return
	}

	_, err := h.UserRepo.GetUserByEmail(req.Username)
	if err == nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthUsernameExists))
		return
	}

	userID, err := h.UserRepo.RegisterUser(req.Username, req.Password, req.Phone)
	if err != nil {
		logger.Error("Failed to register user", zap.Error(err))
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthCreateUserFailed))
		return
	}

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
	var req types.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidRequest))
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthValidationFailed))
		return
	}

	u, err := h.UserRepo.GetUserByEmail(req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidCredentials))
		} else {
			utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthDatabaseError))
		}
		return
	}

	if !encryption.CheckPassword(req.Password, u.Password) {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthInvalidCredentials))
		return
	}

	tokens, err := h.JWTManager.GenerateTokenPair(u.ID)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req types.RefreshRequest
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

	tokens, err := h.JWTManager.GenerateTokenPair(claims.UserID)
	if err != nil {
		utils.RespondError(w, r, h.MsgStore, errors.New(language.AuthTokenGenFailed))
		return
	}

	utils.JSON(w, http.StatusOK, tokens)
}
