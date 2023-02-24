package main

import (
	"magna/routers"
	"magna/services/mangaservice"

	"github.com/robfig/cron"
)

func main() {
	router := routers.InitRouter()

	c := cron.New()
	c.AddFunc("@daily", mangaservice.ClearHotMangaMap)

	router.Run("localhost:8080")
}
