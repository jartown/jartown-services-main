package crawler

import (
	"fmt"

	"github.com/gocolly/colly"
	common "github.com/singkorn/jartown-services-common"
	"github.com/singkorn/jartown-services-main/config"
)

func Crawl() {
	err := common.ConfigInit(&config.Conf)
	if err != nil {
		fmt.Println(err.Error())
	}

	c := colly.NewCollector()
	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	conf := &config.Conf
	fmt.Println(conf.Url)
	c.Visit(conf.Url)
}
