package main

import (
	"magna/routers"
	"magna/services/mangaservice"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
)

// NGUYEN VAN A
// 9704 0000 0000 0018
// 03/07
// OTP

func init() {
	godotenv.Load(".env")
}

func main() {
	// fmt.Println(paymentservice.GetMomoPayURL("pay with MoMo", "http://localhost:8081", "http://localhost:8081", "10000", ""))
	router := routers.InitRouter()

	c := cron.New()
	c.AddFunc("@daily", mangaservice.ClearHotMangaMap)

	port, ok := os.LookupEnv("PORT")

	if !ok {
		port = "8080"
	}

	router.Run(":" + port)
}
