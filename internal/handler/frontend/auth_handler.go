package frontend

import (
	"github.com/gin-gonic/gin"
	"github.com/eduflow/eduflow/internal/pkg/response"
	"github.com/eduflow/eduflow/internal/service"
	"github.com/eduflow/eduflow/pkg/constants"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	token, user, err := h.authService.LoginUser(req.Email, req.Password)
	if err != nil {
		response.Error(c, constants.PasswordErrorCode, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	user, err := h.authService.RegisterUser(req.Email, req.Password, req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	jti, exists := c.Get(constants.ContextKeyJTI)
	if !exists {
		response.Unauthorized(c)
		return
	}

	h.authService.Logout(jti.(string), constants.CacheTTLToken)

	response.Success(c, nil)
}

func (h *AuthHandler) GetDetail(c *gin.Context) {
	userID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c)
		return
	}

	email, _ := c.Get(constants.ContextKeyEmail)

	response.Success(c, gin.H{
		"id":    userID,
		"email": email,
	})
}
