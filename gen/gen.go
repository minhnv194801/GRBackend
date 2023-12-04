package gen

import (
	"log"
	"magna/model"
	"magna/services/ratingservice"
	"magna/services/userservice"
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/goombaio/namegenerator"
)

func GenUsersWithRating(genUserCount, minRatedMangaCount, maxRatedMangaCount int) error {
	mangaList, err := new(model.Manga).GetAllItem()
	if err != nil {
		log.Println(err)
		return err
	}

	for i := 0; i < genUserCount; i++ {
		seed := time.Now().UTC().UnixNano()
		nameGenerator := namegenerator.NewNameGenerator(seed)

		reg, _ := regexp.Compile(`[^\w]`)
		name := nameGenerator.Generate()
		name = reg.ReplaceAllString(name, "")
		id, err := userservice.CreateAccount(name+"@gmail.com", "pass", "Người dùng")
		if err != nil {
			return err
		}
		log.Println(name + "@gmail.com" + " created")
		rand.Seed(time.Now().UnixNano())
		numberOfRatedManga := rand.Intn(maxRatedMangaCount-minRatedMangaCount) + minRatedMangaCount
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(mangaList), func(i, j int) { mangaList[i], mangaList[j] = mangaList[j], mangaList[i] })
		for j := 0; j < numberOfRatedManga; j++ {
			mangaId := mangaList[j].Id.Hex()
			rand.Seed(time.Now().UnixNano())
			rating := rand.Intn(5) + 1
			ratingservice.SetRating(id, mangaId, rating)
			log.Println(name + "@gmail.com" + " rated " + mangaList[j].Name + " " + strconv.Itoa(rating) + " star")
		}
	}

	return nil
}
