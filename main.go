package main

import (
	"magna/routers"
)

func main() {
	// router := gin.Default()

	// corsConfig := cors.DefaultConfig()
	// corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	// corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	// // corsConfig.AllowHeaders = []string{"Content-Type", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"}
	// // To be able to send tokens to the server.
	// corsConfig.AllowCredentials = true
	// // OPTIONS method for ReactJS
	// corsConfig.AddAllowMethods("OPTIONS")
	// // Register the middleware
	// router.Use(cors.New(corsConfig))

	// router.Group("/api/v1")
	// router.GET("/chapterlist/:mangaid", routes.GetChapterInfo)
	// router.GET("/chapter/:chapterid", routes.GetChapterList)

	router := routers.InitRouter()
	router.Run("localhost:8080")
	// db, _ := database.GetMongoDB()
	// coll := db.Collection("Manga")

	// cursor, err := coll.Find(context.TODO(), bson.D{{}})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer cursor.Close(context.Background())
	// for cursor.Next(context.Background()) {
	// 	// To decode into a struct, use cursor.Decode()
	// 	var result model.Manga
	// 	err := cursor.Decode(&result)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	var doc model.Chapter
	// 	fmt.Println(result.Chapters[len(result.Chapters)-1])
	// 	filter := bson.D{{"_id", result.Chapters[len(result.Chapters)-1]}}
	// 	opts := options.FindOne()
	// 	chapterColl := db.Collection("Chapter")
	// 	found := chapterColl.FindOne(context.TODO(), filter, opts)
	// 	if found.Err() != nil {
	// 		log.Println("Hello")
	// 		log.Fatal(found.Err())
	// 	}
	// 	err = found.Decode(&doc)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	result.UpdateTime = doc.UpdateTime
	// 	result.UpdateUpdateTime()
	// 	fmt.Println(result.Id)
	// }
}
