package main

import (
	"log"
	"rss-reader/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

var conf *viper.Viper = viper.New()

func main() {
	if err := config.LoadConfig(conf, "PORT", "LOGIN", "PASSWORD"); err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(middleware.BasicAuth(validateBasicAuth))
	e.POST("/login", func(c echo.Context) error {
		return c.HTML(200, "OK")
	})
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, "hello")
	})

	if err := e.Start(":" + conf.GetString("PORT")); err != nil {
		log.Printf("Error in echo: %v", err)
	}
}

func validateBasicAuth(user, pass string, c echo.Context) (bool, error) {
	if user == conf.GetString("LOGIN") && pass == conf.GetString("PASSWORD") {
		return true, nil
	}
	return false, nil
}
