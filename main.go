package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

func main() {
	SetupViper()
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.HTML(200, "hello")
	})

	err := e.Start(":" + viper.GetString("PORT"))
	log.Printf("Error in echo: %v", err)
}

func SetupViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to setup viper: %v", err)
	}
	requiredKeys := []string{"PORT"}
	for _, key := range requiredKeys {
		if !viper.IsSet(key) {
			log.Fatalf("Key '%s' is not set", key)
		}
	}
}
