package main

import (
	"log"
	"magna/routers"
	"magna/services/mangaservice"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
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
