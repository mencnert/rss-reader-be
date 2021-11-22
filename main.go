package main

import (
	"log"
	"rss-reader/config"
	"time"

	cron "github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

var conf *viper.Viper = viper.New()

func main() {
	if err := config.LoadConfig(conf, "PORT", "LOGIN", "PASSWORD", "RSS_FETCH_EVERY_N_SEC"); err != nil {
		log.Fatalf("Unable to load config: %v", err)
	}

	if err := startRssFetchCronJob(conf.GetInt("RSS_FETCH_EVERY_N_SEC")); err != nil {
		log.Fatalf("Error to setup rss fetch cron job: %v", err)
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

func fetchRss() {
	log.Println("TODO fetch new rss")
}

func startRssFetchCronJob(rssFetchEveryNSecs int) error {
	sch := cron.NewScheduler(time.UTC)
	_, err := sch.Every(rssFetchEveryNSecs).Seconds().Do(fetchRss)
	if err != nil {
		return err
	}
	sch.StartAsync()
	return nil
}
