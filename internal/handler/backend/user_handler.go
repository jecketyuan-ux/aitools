package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/eduflow/eduflow/internal/domain"
	"github.com/eduflow/eduflow/internal/pkg/crypto"
	"github.com/eduflow/eduflow/internal/pkg/response"
	"github.com/eduflow/eduflow/internal/service"
	"github.com/eduflow/eduflow/pkg/utils"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	page, size := utils.GetPageParams(c.Query("page"), c.Query("size"))

	filters := make(map[string]interface{})
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}
	if email := c.Query("email"); email != "" {
		filters["email"] = email
	}

	users, total, err := h.userService.List(page, size, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPagination(c, users, total, page, size)
}

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	hashedPassword, err := crypto.HashPassword(req.Password, salt)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	user := &domain.User{
		Email:         req.Email,
		Name:          req.Name,
		Password:      hashedPassword,
		Salt:          salt,
		IsActive:      1,
		IsSetPassword: 1,
	}

	if err := h.userService.Create(user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.Success(c, user)
}

type UpdateUserRequest struct {
	Name   string `json:"name" binding:"required"`
	Avatar *int   `json:"avatar"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	user.Name = req.Name
	user.Avatar = req.Avatar

	if err := h.userService.Update(user); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	if err := h.userService.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
