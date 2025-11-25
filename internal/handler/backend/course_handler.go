package backend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/playedu/playedu-go/internal/domain"
	"github.com/playedu/playedu-go/internal/pkg/response"
	"github.com/playedu/playedu-go/internal/service"
	"github.com/playedu/playedu-go/pkg/utils"
)

type CourseHandler struct {
	courseService service.CourseService
}

func NewCourseHandler(courseService service.CourseService) *CourseHandler {
	return &CourseHandler{
		courseService: courseService,
	}
}

func (h *CourseHandler) List(c *gin.Context) {
	page, size := utils.GetPageParams(c.Query("page"), c.Query("size"))

	filters := make(map[string]interface{})
	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}
	if isShow := c.Query("is_show"); isShow != "" {
		if val, err := strconv.Atoi(isShow); err == nil {
			filters["is_show"] = val
		}
	}

	courses, total, err := h.courseService.List(page, size, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPagination(c, courses, total, page, size)
}

type CreateCourseRequest struct {
	Title      string `json:"title" binding:"required"`
	ShortDesc  string `json:"short_desc"`
	Thumb      *int   `json:"thumb"`
	IsRequired int    `json:"is_required"`
	IsShow     int    `json:"is_show"`
}

func (h *CourseHandler) Create(c *gin.Context) {
	var req CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	course := &domain.Course{
		Title:      req.Title,
		ShortDesc:  req.ShortDesc,
		Thumb:      req.Thumb,
		IsRequired: req.IsRequired,
		IsShow:     req.IsShow,
	}

	if err := h.courseService.Create(course); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, course)
}

func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	course, err := h.courseService.GetByID(id)
	if err != nil {
		response.NotFound(c, "course not found")
		return
	}

	response.Success(c, course)
}

type UpdateCourseRequest struct {
	Title      string `json:"title" binding:"required"`
	ShortDesc  string `json:"short_desc"`
	Thumb      *int   `json:"thumb"`
	IsRequired int    `json:"is_required"`
	IsShow     int    `json:"is_show"`
}

func (h *CourseHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	course, err := h.courseService.GetByID(id)
	if err != nil {
		response.NotFound(c, "course not found")
		return
	}

	var req UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, err.Error())
		return
	}

	course.Title = req.Title
	course.ShortDesc = req.ShortDesc
	course.Thumb = req.Thumb
	course.IsRequired = req.IsRequired
	course.IsShow = req.IsShow

	if err := h.courseService.Update(course); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, course)
}

func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.ParamError(c, "invalid id")
		return
	}

	if err := h.courseService.Delete(id); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
