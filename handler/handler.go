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
	router *gin.Engine
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

	// shopRG := apiV1RG.Group("/shop")
	// h.registerShopSvc(shopRG)

	// orderRG := apiV1RG.Group("/order")
	// h.registerOrderSvc(orderRG)

	// websocketV1RG := apiV1RG.Group("/ws")
	// h.registerWebsocket(websocketV1RG)

	common.HandlerRegisterSwagger(h.router)

	return h
}
