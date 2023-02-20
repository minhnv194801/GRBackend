package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chapter struct {
	Id         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Manga      primitive.ObjectID   `bson:"manga"`
	Name       string               `bson:"name" json:"title"`
	Cover      string               `bson:"cover"`
	Price      uint                 `bson:"price"`
	UpdateTime uint                 `bson:"updateTime"`
	Images     []string             `bson:"images"`
	OwnedUsers []primitive.ObjectID `bson:"ownedUsers"`
	Reports    []primitive.ObjectID `bson:"reports"`
}

func (chapter *Chapter) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return [12]byte{}, err
	}

	id, err := getExistedChapterID(chapter.Manga, chapter.Name)
	if id != primitive.NilObjectID {
		chapter.Id = id
		return id, err
	}

	result, err := coll.InsertOne(context.TODO(), chapter)
	if err != nil {
		return [12]byte{}, err
	}

	chapter.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (chapter *Chapter) GetFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(chapter)
	if err != nil {
		return err
	}

	return nil
}

func (chapter *Chapter) GetItemListFromObjectId(objID []primitive.ObjectID) ([]Chapter, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Chapter, 0)
	multiFilter := bson.M{"_id": bson.M{"$in": objID}}
	multiFindOpts := options.Find().SetSort(bson.D{{"updateTime", 1}})
	cursor, err := coll.Find(context.TODO(), multiFilter, multiFindOpts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	return listItem, nil
}

func getExistedChapterID(mangaId primitive.ObjectID, name string) (primitive.ObjectID, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return [12]byte{}, err
	}

	var doc Chapter
	filter := bson.D{{"manga", mangaId}, {"name", name}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return [12]byte{}, found.Err()
	}
	err = found.Decode(&doc)
	if err != nil {
		return [12]byte{}, err
	}
	return doc.Id, nil
}
