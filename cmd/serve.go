package cmd

import (
	"fmt"
	"log"
	"net/http"
	repo "rss-reader/repositories"
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
	rssRepo      repo.RssRepository
)

func newWebCmd() *cobra.Command {
	webCmd := &cobra.Command{
		Use: "web",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := loadConfig(webCfg); err != nil {
				return err
			}
			if err := loadConfig(dbCfg); err != nil {
				return err
			}
			if err := loadConfig(schedulerCfg); err != nil {
				return err
			}
			if err := startRssFetchCronJob(schedulerCfg.RssFetchEveryNSecs); err != nil {
				return err
			}
			if err := startCleanDbCronJob(schedulerCfg.CleanDbEveryNHours); err != nil {
				return err
			}

			rssRepo = repo.NewPgRssRepository(dbCfg.DBUrl)
			return rssRepo.Open()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			e := echo.New()
			e.Use(newCorsMiddlewareWithOrigins(webCfg.CorsAllowOrigins))
			e.Use(middleware.BasicAuth(validateBasicAuth))
			e.POST("/checkauth", httpCheckAuth)
			e.GET("/rss", httpGetRss)
			e.PUT("/rss/:rss_id", httpUpdateRss)

			return e.Start(fmt.Sprintf(":%d", webCfg.Port))
		},
		PostRunE: func(cmd *cobra.Command, args []string) error {
			return rssRepo.Close()
		},
	}
	return webCmd
}

func newCorsMiddlewareWithOrigins(origins []string) echo.MiddlewareFunc {
	corsConfig := middleware.DefaultCORSConfig
	corsConfig.AllowOrigins = webCfg.CorsAllowOrigins
	return middleware.CORSWithConfig(corsConfig)
}

func validateBasicAuth(user, pass string, c echo.Context) (bool, error) {
	if user == webCfg.Username && pass == webCfg.Password {
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
	log.Println("Starting rss fetch scheduler")
	sch := cron.NewScheduler(time.UTC)
	if _, err := sch.Every(rssFetchEveryNSecs).Seconds().Do(fetchRss); err != nil {
		return err
	}
	sch.StartAsync()
	return nil
}

func startCleanDbCronJob(cleanDbEveryNHours int) error {
	log.Println("Starting clean DB scheduler")
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

func httpUpdateRss(c echo.Context) error {
	i := c.Param("rss_id")
	rssId, err := strconv.Atoi(i)
	if err != nil {
		c.HTML(http.StatusBadRequest, fmt.Sprintf("Invalid path parameter rss_id of type int: '%s'", i))
		return nil
	}
	var rssDto repo.RssDTO
	if err := c.Bind(&rssDto); err != nil {
		return err
	}
	rssDto.Id = rssId
	if err := rssRepo.Update(rssDto); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, struct{ Status string }{"OK"})
}
