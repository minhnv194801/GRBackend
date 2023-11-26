package main

import (
	"magna/routers"
	"magna/services/mangaservice"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	router := routers.InitRouter()

	c := cron.New()
	c.AddFunc("@daily", mangaservice.ClearHotMangaMap)

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}

	router.Run(":" + port)
}
