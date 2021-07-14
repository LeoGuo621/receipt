package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"receipt/routes"
)


// global Echo instance
var Pecho *echo.Echo
func main() {
	Pecho = echo.New()
	// install logger middleware
	Pecho.Use(middleware.Logger())
	// install recover middleware
	Pecho.Use(middleware.Recover())
	// install gzip middleware
	Pecho.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	// frontend static file service
	// staticFile()

	Pecho.GET("/ping", routes.PingHandler)
	Pecho.POST("/register", routes.Register)
	Pecho.POST("/login", routes.Login)
	Pecho.GET("/session", routes.Session)

	Pecho.POST("/content", routes.Upload)
	Pecho.GET("/content", routes.GetContents)
	// TODO: deploys contract


	Pecho.Logger.Fatal(Pecho.Start(":1323"))
}
