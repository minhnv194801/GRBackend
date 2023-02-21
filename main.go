package main

import (
	"magna/routers"
)

func main() {
	router := routers.InitRouter()
	router.Run("localhost:8080")
}
