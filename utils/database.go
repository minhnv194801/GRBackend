package utils

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
)

func FixCorruptedImages(oldImage, newImage string) error {
	mangaColl, err := database.GetMangaCollection()
	if err != nil {
		return err
	}
	chapterColl, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"cover", oldImage}}
	update := bson.D{{"$set", bson.D{{"cover", newImage}}}}
	_, err = mangaColl.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	_, err = chapterColl.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}
