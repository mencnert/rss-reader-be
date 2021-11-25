package cmd

import (
	"fmt"
	"log"
	"time"

	cron "github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

func newWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use: "web",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Starting schedulers")
			if err := startRssFetchCronJob(Config.RssFetchEveryNSecs); err != nil {
				log.Printf("Error to setup rss fetch cron job: %v", err)
				return err
			}
			if err := startCleanDbCronJob(Config.CleanDbEveryNHours); err != nil {
				log.Printf("Error to setup clean db cron job: %v", err)
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			e := echo.New()
			e.Use(middleware.CORS())
			e.Use(middleware.BasicAuth(validateBasicAuth))
			e.POST("/checkauth", httpCheckAuth)
			e.GET("/", httpSth)

			return e.Start(fmt.Sprintf(":%d", Config.Port))
		},
	}
	return webCmd
}

func validateBasicAuth(user, pass string, c echo.Context) (bool, error) {
	if user == Config.Username && pass == Config.Password {
		return true, nil
	}
	return false, nil
}

func fetchRss() {
	log.Println("TODO fetch new rss")
}

func cleanDb() {
	log.Println("TODO clean db")
}

func startRssFetchCronJob(rssFetchEveryNSecs int) error {
	sch := cron.NewScheduler(time.UTC)
	if _, err := sch.Every(rssFetchEveryNSecs).Seconds().Do(fetchRss); err != nil {
		return err
	}
	sch.StartAsync()
	return nil
}

func startCleanDbCronJob(cleanDbEveryNHours int) error {
	sch := cron.NewScheduler(time.UTC)
	if _, err := sch.Every(cleanDbEveryNHours).Hours().Do(cleanDb); err != nil {
		return err
	}
	go startWithDeplay(sch, 30*time.Second)
	return nil
}

func startWithDeplay(scheduler *cron.Scheduler, delay time.Duration) {
	time.Sleep(delay)
	scheduler.StartAsync()
}

func httpCheckAuth(c echo.Context) error {
	return c.HTML(200, "OK")
}

func httpSth(c echo.Context) error {
	return c.HTML(200, "hello")
}
