package main

import (
	"github.com/gin-gonic/gin"
	"gitlab.ssc.dev/Research/tia/code/api"
	"gitlab.ssc.dev/Research/tia/code/db"
)

func main() {
	db.Main()
	r := gin.Default()
	r.NoRoute(api.FOF)

	//Test route
	r.GET("/ping", api.Ping)
	//Static Routes
	r.GET("/", api.Home)
	r.GET("/home", api.Home)
	r.GET("/about", api.About)
	r.GET("/help", api.Help)
	//Dynamic Routes
	r.GET("/graph", api.Graph)
	r.GET("/graph/:vendor", api.GetProducts)
	r.GET("/graph/:vendor/:product", api.GetVulns)
	r.GET("/graph/:vendor/:product/stix", api.GetStix)
	r.GET("/graph/:vendor/:product/scores", api.GetScores)
	r.GET("/graph/:vendor/:product/cwes", api.GetCwes)

	//Load HTML Templates
	r.Static("/assets", "../website/assets")
	r.LoadHTMLGlob("../website/templates/**/*.html")
	// r.LoadHTMLGlob("../stixviz/*.html")

	//
	r.GET("/test", api.TestViz)

	r.Run()
}

func Hello() string {
	return "Hello, world."
}
