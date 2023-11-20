package userservice

import (
	"errors"
	"log"
	"magna/model"
	"magna/utils"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/cdipaolo/goml/cluster"
	"github.com/wilcosheh/tfidf"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserInfo(userId string) (*model.User, error) {
	objId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}
	user := new(model.User)
	err = user.GetItemFromObjectId(objId)
	if err != nil {
		return nil, err
	}
	return user, err
}

func CreateAccount(email, password, role string) (string, error) {
	if !utils.ValidateEmail(email) {
		return "", errors.New("Email không hợp lệ")
	}
	if !utils.ValidatePassword(password) {
		return "", errors.New("Password không hợp lệ")
	}

	user := new(model.User)
	_, err := user.CreateNewUser(email, password, role)
	if err != nil {
		log.Println(err.Error(), "err.Error() services/userservice/user_service.go:38")
		return "", err
	}

	id := user.Id.Hex()

	return id, nil
}

func GetUserRecommendations(userId string, count int) ([]model.Manga, error) {
	if userId == "" {
		return GetRandomUserRecommendation(count)
	}

	var k = 3

	mangaList, err := new(model.Manga).GetAllItem()
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var wordMap = make(map[string]int)
	f := tfidf.New()
	for _, manga := range mangaList {
		for index, tag := range manga.Tags {
			manga.Tags[index] = strings.Replace(tag, " ", "", -1)
			_, ok := wordMap[manga.Tags[index]]
			if !ok {
				wordMap[manga.Tags[index]] = len(wordMap)
			}
		}
		f.AddDocs(strings.Join(manga.Tags, " "))
	}

	user := new(model.User)
	userIdString, _ := primitive.ObjectIDFromHex(userId)
	err = user.GetItemFromObjectId(userIdString)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	var trainset [][]float64
	var expectedResults []float64
	for mangaId, rate := range user.Rate {
		manga := new(model.Manga)
		err = manga.GetItemFromObjectId(mangaId)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		for index, tag := range manga.Tags {
			manga.Tags[index] = strings.Replace(tag, " ", "", -1)
		}
		features := f.Cal(strings.Join(manga.Tags, " "))
		weights := make([]float64, len(wordMap))
		for word, weight := range features {
			weights[wordMap[word]] = weight
		}

		trainset = append(trainset, weights)
		expectedResults = append(expectedResults, float64(rate))
	}
	knn := cluster.NewKNN(k, trainset, expectedResults, utils.CosineDistance)

	predictedValue := make(map[primitive.ObjectID]float64)
	for _, manga := range mangaList {
		_, ok := user.Rate[manga.Id]
		if ok {
			continue
		}
		features := f.Cal(strings.Join(manga.Tags, " "))
		weights := make([]float64, len(wordMap))
		for word, weight := range features {
			weights[wordMap[word]] = weight
		}

		predicted, err := knn.Predict(weights)
		if err != nil {
			log.Println(err)
			mangaList, _ = GetRandomUserRecommendation(count)
			return mangaList, err
		}

		predictedValue[manga.Id] = predicted[0]
	}

	keys := make([]primitive.ObjectID, 0, len(predictedValue))
	for key := range predictedValue {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return predictedValue[keys[i]] > predictedValue[keys[j]] })

	recommendedIds := keys[:utils.Min(len(keys), count)]
	var recommendedItems []model.Manga
	for _, id := range recommendedIds {
		item := new(model.Manga)
		item.GetItemFromObjectId(id)
		recommendedItems = append(recommendedItems, *item)
	}

	return recommendedItems, nil
}

func GetRandomUserRecommendation(count int) ([]model.Manga, error) {
	mangaList, err := new(model.Manga).GetAllItem()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(mangaList), func(i, j int) { mangaList[i], mangaList[j] = mangaList[j], mangaList[i] })
	return mangaList[:utils.Min(len(mangaList), count)], nil
}

func GetTotalCount() (int, error) {
	return new(model.User).GetTotalCount()
}
