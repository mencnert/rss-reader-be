package cmd

import (
	"fmt"
	"log"
	"net/http"
	"rss-reader/rss"
	"strconv"
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
			e.GET("/rss", httpGetRss)
			e.PUT("/rss/:id", httpChangeViewedState)

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
		rssRepo.SaveOrUpdateAll(feed.Entries)
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

func httpGetRss(c echo.Context) error {
	data, err := rssRepo.GetAll()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, data)
}

func httpChangeViewedState(c echo.Context) error {
	i := c.Param("id")
	id, err := strconv.Atoi(i)
	if err != nil {
		c.HTML(http.StatusBadRequest, fmt.Sprintf("Invalid path parameter id of type int: '%s'", i))
		return nil
	}

	v := c.QueryParam("viewed")
	viewed, err := strconv.ParseBool(v)
	if err != nil {
		c.HTML(http.StatusBadRequest, fmt.Sprintf("Invalid query parameter 'viewed' of type bool: '%s'", v))
		return nil
	}
	if err := rssRepo.UpdateViewedById(id, viewed); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct{ Status string }{"OK"})
}
