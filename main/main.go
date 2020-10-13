package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	common "github.com/singkorn/jartown-services-common"
	"github.com/singkorn/jartown-services-common/log"
	"github.com/singkorn/jartown-services-main/config"
	"github.com/singkorn/jartown-services-main/crawler"
	"github.com/singkorn/jartown-services-main/database"
)

func main() {
	crawler.Crawl()
	err := common.ConfigInit(&config.Conf)
	if err != nil {
		log.Fatal(err)
	}
	log.Init(config.Conf.Env, config.ApplicationName, "error")

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(log.GinLogger())

	database.Init()

	// pubsub := pubsub.InitBus()

	// h := handler.Init(r)
	// h.ShopSvc = service.ShopServiceInit(database.GetShopDatabase())
	// h.OrderSvc = service.OrderServiceInit(database.GetOrderDatabase(), pubsub)

	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		log.Fatal(err)
	}
}
