package handler

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	common "github.com/singkorn/jartown-services-common"
	"github.com/singkorn/jartown-services-common/auth"
	// "github.com/singkorn/jartown-services-main/service"
)

type Route struct {
	Name        string
	Description string
	Method      string
	Pattern     string
	Handler     gin.HandlerFunc
}

type Handler struct {
	router        *gin.Engine
	ListReviewSvc service.ListReviewService
	// ShopSvc  service.ShopService
	// OrderSvc service.OrderService
}

func Init(r *gin.Engine) *Handler {
	h := &Handler{}
	h.router = r

	corsConf := cors.DefaultConfig()
	corsConf.AddAllowHeaders("Authorization")
	corsConf.AllowAllOrigins = true // TODO: Stricter CORS
	h.router.Use(cors.New(corsConf))

	h.router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, common.ErrorReply{
			ErrorMsg: common.ErrEndpointNotFound.Error(),
		})
	})

	apiV1RG := h.router.Group("/api/v1")
	apiV1RG.Use(auth.MiddlewareAuthentication)
	apiV1RG.Use(auth.MiddlewareAuthRequired)

	listReviewRG := h.router.Group("/listreview")
	h.registerListReviewSvc(listReviewRG)

	common.HandlerRegisterSwagger(h.router)

	return h
}
