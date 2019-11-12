package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	common "github.com/singkorn/jartown-services-common"
)

func (h *Handler) registerListReviewSvc(r *gin.RouterGroup) {
	for _, route := range h.listReviewSvcRoutes() {
		r.Handle(route.Method, route.Pattern, route.Handler)
	}
}

func (h *Handler) listReviewSvcRoutes() map[string]Route {
	return map[string]Route{
		"GetListReview": {
			Name:        "GetListReview",
			Description: "Get list and review data by ID",
			Method:      http.MethodGet,
			Pattern:     "/:id",
			Handler:     h.getListReview,
		},
	}
}

func (h *Handler) getListReview(c *gin.Context) {
	doer, _ := c.Get("JartownUser")
	id := c.Param("id")
	listReview, err := h.ListReviewSvc.GetListReview(doer, id)
	if err != nil {
		common.HandlerReturnError(c, err)
	} else {
		common.HandlerReturnData(c, listReview)
	}
}
