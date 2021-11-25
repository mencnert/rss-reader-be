package cmd

import (
	"fmt"
	"log"
	"rss-reader/rss"
	"time"

	cron "github.com/go-co-op/gocron"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

var (
	rssConnector = rss.Connector{}
	feeds        = []string{"https://stackoverflow.com/feeds/tag?tagnames=go&sort=newest"}
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
			e.Use(newCorsMiddlewareWithOrigins(Config.CorsAllowOrigins))
			e.Use(middleware.BasicAuth(validateBasicAuth))
			e.POST("/checkauth", httpCheckAuth)
			e.GET("/", httpSth)

			return e.Start(fmt.Sprintf(":%d", Config.Port))
		},
	}
	return webCmd
}

func newCorsMiddlewareWithOrigins(origins []string) echo.MiddlewareFunc {
	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowOrigins = Config.CorsAllowOrigins
	return middleware.CORSWithConfig(corsConfig)
}

func validateBasicAuth(user, pass string, c echo.Context) (bool, error) {
	if user == Config.Username && pass == Config.Password {
		return true, nil
	}
	return false, nil
}

func fetchRss() {
	log.Println("Updating rss")
	for _, url := range feeds {
		feed, err := rssConnector.Fetch(url)
		if err != nil {
			log.Printf("Error during fetch: %v\n", err)
		}
		//TODO: move to separate package
		stmt, err := DB.Prepare("INSERT INTO rss (url, rank, title, last_update) VALUES ($1, $2, $3, NOW()) ON CONFLICT (url) DO UPDATE SET rank = excluded.rank, title = excluded.title, last_update = NOW();")
		if err != nil {
			log.Printf("Error during preparation of query: %v\n", err)
			return
		}
		for _, entry := range feed.Entries {
			if _, err := stmt.Exec(entry.Url, entry.Rank, entry.Title); err != nil {
				log.Printf("Error during execution of query: %v\n", err)
				return
			}
		}
	}
	log.Println("Updating rss done")
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
	go startWithDelay(sch, 30*time.Second)
	return nil
}

func startWithDelay(scheduler *cron.Scheduler, delay time.Duration) {
	time.Sleep(delay)
	scheduler.StartAsync()
}

func httpCheckAuth(c echo.Context) error {
	return c.HTML(200, "OK")
}

func httpSth(c echo.Context) error {
	return c.HTML(200, "hello")
}
