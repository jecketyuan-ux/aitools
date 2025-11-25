package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/playedu/playedu-go/internal/pkg/response"
	"github.com/playedu/playedu-go/internal/service"
	"github.com/playedu/playedu-go/pkg/constants"
	"github.com/playedu/playedu-go/pkg/utils"
)

type ResourceHandler struct {
	resourceService service.ResourceService
}

func NewResourceHandler(resourceService service.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

func (h *ResourceHandler) List(c *gin.Context) {
	page, size := utils.GetPageParams(c.Query("page"), c.Query("size"))

	filters := make(map[string]interface{})
	if resourceType := c.Query("type"); resourceType != "" {
		filters["type"] = resourceType
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		if val, err := strconv.Atoi(categoryID); err == nil {
			filters["category_id"] = val
		}
	}
	if name := c.Query("name"); name != "" {
		filters["name"] = name
	}

	resources, total, err := h.resourceService.List(page, size, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPagination(c, resources, total, page, size)
}

func (h *ResourceHandler) UploadVideo(c *gin.Context) {
	adminID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.ParamError(c, "file is required")
		return
	}

	resource, err := h.resourceService.UploadVideo(file, adminID.(int))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resource)
}

func (h *ResourceHandler) UploadImage(c *gin.Context) {
	adminID, exists := c.Get(constants.ContextKeyUserID)
	if !exists {
		response.Unauthorized(c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.ParamError(c, "file is required")
		return
	}

	resource, err := h.resourceService.UploadImage(file, adminID.(int))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, resource)
}

func (h *ResourceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	if err := h.resourceService.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
