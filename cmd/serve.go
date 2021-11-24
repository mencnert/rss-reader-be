package cmd

import (
	"log"
	"time"

	cron "github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use: "web",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := startRssFetchCronJob(viper.GetInt("RSS_FETCH_EVERY_N_SECS")); err != nil {
				log.Printf("Error to setup rss fetch cron job: %v", err)
				return err
			}
			if err := startCleanDbCronJob(viper.GetInt("CLEAN_DB_EVERY_N_HOURS")); err != nil {
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

			return e.Start(":" + viper.GetString("PORT"))
		},
	}
	return webCmd
}

func validateBasicAuth(user, pass string, c echo.Context) (bool, error) {
	if user == viper.GetString("LOGIN") && pass == viper.GetString("PASSWORD") {
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
	go startWithDeplay(sch, 5000)
	return nil
}

func startWithDeplay(scheduler *cron.Scheduler, msDelay int) {
	time.Sleep(time.Duration(msDelay) * time.Millisecond)
	scheduler.StartAsync()
}

func httpCheckAuth(c echo.Context) error {
	return c.HTML(200, "OK")
}

func httpSth(c echo.Context) error {
	return c.HTML(200, "hello")
}
