package model

import (
	"context"
	"magna/database"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Manga struct {
	Id            primitive.ObjectID   `bson:"_id,omitempty"`
	Name          string               `bson:"name"`
	AlternateName []string             `bson:"alternateName"`
	Author        []string             `bson:"author"`
	Cover         string               `bson:"cover"`
	Description   string               `bson:"description"`
	Status        Status               `bson:"status"`
	UpdateTime    uint                 `bson:"updateTime"`
	IsRecommended bool                 `bson:"isRecommended"`
	Tags          []string             `bson:"tags"`
	FollowedUsers []primitive.ObjectID `bson:"followedUsers"`
	Chapters      []primitive.ObjectID `bson:"chapters"`
	Comments      []primitive.ObjectID `bson:"comments"`
}

type Status int

const (
	Ongoing Status = iota
	Finished
)

func (manga *Manga) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return [12]byte{}, err
	}

	id, err := getExistedTitleID(manga.Name)
	if id != primitive.NilObjectID {
		manga.Id = id
		return id, err
	}

	result, err := coll.InsertOne(context.TODO(), manga)
	if err != nil {
		return [12]byte{}, err
	}

	manga.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (manga *Manga) UpdateChapter(chapter *Chapter) error {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return err
	}

	manga.Chapters = append(manga.Chapters, chapter.Id)
	if manga.UpdateTime < chapter.UpdateTime {
		manga.UpdateTime = chapter.UpdateTime
	}
	filter := bson.D{{"_id", manga.Id}}
	update := bson.D{{"$set", bson.D{{"chapters", manga.Chapters}, {"updateTime", manga.UpdateTime}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (manga *Manga) GetItemFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return err
	}
	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(&manga)
	if err != nil {
		return err
	}
	return nil
}

func (manga *Manga) GetItemFromName(name string) error {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "name", Value: name}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(&manga)
	if err != nil {
		return err
	}
	return nil
}

func (manga *Manga) GetItemListFromObjectId(objID []primitive.ObjectID) ([]Manga, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Manga, 0)
	aggregatePipeline := bson.A{}
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$match",
				bson.M{"_id": bson.M{"$in": objID}},
			},
		})
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$addFields",
				bson.D{
					{"order",
						bson.D{
							{"$indexOfArray",
								bson.A{
									objID,
									"$_id",
								},
							},
						},
					},
				},
			},
		})
	aggregatePipeline = append(aggregatePipeline, bson.D{{"$sort", bson.D{{"order", 1}}}})
	cursor, err := coll.Aggregate(context.TODO(), aggregatePipeline)
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

func (manga *Manga) GetRandomExcludedItemListFromObjectId(objID []primitive.ObjectID, count int) ([]Manga, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Manga, 0)
	aggregatePipeline := bson.A{}
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$match",
				bson.M{
					"_id": bson.M{
						"$nin": objID,
					},
				},
			},
		})
	aggregatePipeline = append(aggregatePipeline, bson.D{{"$sample", bson.D{{"size", 5}}}})
	cursor, err := coll.Aggregate(context.TODO(), aggregatePipeline)
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

func (manga *Manga) GetListRecommendManga(count int) ([]Manga, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Manga, 0)
	filter := bson.M{"isRecommended": true}
	cursor, err := coll.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	if count < len(listItem) {
		return listItem[:count], nil
	} else {
		return listItem[:], nil
	}
}

func (manga *Manga) GetNewestItemList(count int) ([]Manga, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Manga, 0)
	filter := bson.D{{}}
	opts := options.Find().SetSort(bson.D{{"updateTime", -1}})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	if count < len(listItem) {
		return listItem[:count], nil
	} else {
		return listItem[:], nil
	}
}

func getExistedTitleID(name string) (primitive.ObjectID, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return [12]byte{}, err
	}

	var doc Manga
	filter := bson.D{primitive.E{Key: "name", Value: name}}
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
