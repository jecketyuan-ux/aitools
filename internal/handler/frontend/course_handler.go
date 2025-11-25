package frontend

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/eduflow/eduflow/internal/pkg/response"
	"github.com/eduflow/eduflow/internal/service"
	"github.com/eduflow/eduflow/pkg/utils"
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
	filters["is_show"] = 1

	if title := c.Query("title"); title != "" {
		filters["title"] = title
	}

	courses, total, err := h.courseService.List(page, size, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.SuccessPagination(c, courses, total, page, size)
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

	if course.IsShow != 1 {
		response.NotFound(c, "course not found")
		return
	}

	response.Success(c, course)
}
