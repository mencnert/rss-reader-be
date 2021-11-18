package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

var (
	// port = GetEnvOrPanic("PORT")
	port = "5000"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, "hello")
	})

	e.Start(":" + port)
}

func GetEnvOrPanic(name string) string {
	env := os.Getenv(name)
	if env == "" {
		log.Fatalf("Env '%s' is empty", name)
	}
	return env
}
