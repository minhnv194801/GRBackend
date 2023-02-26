package model

import (
	"context"
	"magna/database"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Manga struct {
	Id            primitive.ObjectID         `bson:"_id,omitempty"`
	Name          string                     `bson:"name"`
	AlternateName []string                   `bson:"alternateName"`
	Author        []string                   `bson:"author"`
	Cover         string                     `bson:"cover"`
	Description   string                     `bson:"description"`
	Status        Status                     `bson:"status"`
	UpdateTime    uint                       `bson:"updateTime"`
	IsRecommended bool                       `bson:"isRecommended"`
	Tags          []string                   `bson:"tags"`
	FollowedUsers []primitive.ObjectID       `bson:"followedUsers"`
	Chapters      []primitive.ObjectID       `bson:"chapters"`
	Comments      []primitive.ObjectID       `bson:"comments"`
	Rated         map[primitive.ObjectID]int `bson:"rated"`
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

func (manga *Manga) SetRated() error {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", manga.Id}}
	update := bson.D{{"$set", bson.D{{"rated", manga.Rated}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (manga *Manga) SetUserFavorite(user primitive.ObjectID) error {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return err
	}
	for index, followedUser := range manga.FollowedUsers {
		if followedUser == user {
			ret := make([]primitive.ObjectID, 0)
			ret = append(ret, manga.FollowedUsers[:index]...)
			manga.FollowedUsers = append(ret, manga.FollowedUsers[index+1:]...)
			filter := bson.D{{"_id", manga.Id}}
			update := bson.D{{"$set", bson.D{{"followedUsers", manga.FollowedUsers}}}}
			_, err = coll.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}
			return nil
		}
	}

	manga.FollowedUsers = append(manga.FollowedUsers, user)
	filter := bson.D{{"_id", manga.Id}}
	update := bson.D{{"$set", bson.D{{"followedUsers", manga.FollowedUsers}}}}
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
	aggregatePipeline = append(aggregatePipeline, bson.D{{"$sample", bson.D{{"size", count}}}})
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

func (manga *Manga) GetNewestItemList(position, count int) ([]Manga, int, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Manga, 0)
	filter := bson.D{{}}
	opts := options.Find()
	opts.SetSort(bson.M{"updateTime": -1})
	opts.SetSkip(int64(position))
	opts.SetLimit(int64(count))

	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if count < len(listItem) {
		return listItem[:count], int(totalCount), nil
	} else {
		return listItem[:], int(totalCount), nil
	}
}

func (manga *Manga) GetTotalCount() (int, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return 0, err
	}

	filter := bson.D{{}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (manga *Manga) Filter(query string, tags []string, position, count int) ([]Manga, int, error) {
	coll, err := database.GetMangaCollection()
	if err != nil {
		return nil, 0, err
	}

	filter := bson.M{}
	if len(tags) != 0 {
		filter["tags"] = bson.D{{"$all", tags}}
	}
	if strings.Trim(query, " ") != "" {
		filter["$or"] = []interface{}{
			bson.D{{"name", primitive.Regex{Pattern: query, Options: "i"}}},
			bson.D{{"description", primitive.Regex{Pattern: query, Options: "i"}}},
		}
	}
	opts := options.Find().SetSkip(int64(position)).SetLimit(int64(count)).SetSort(bson.M{"updateTime": -1})
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.Background())
	var listItem []Manga
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, 0, err
	}

	totalCount, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return nil, 0, err
	}

	if count < len(listItem) {
		return listItem[:count], int(totalCount), nil
	} else {
		return listItem[:], int(totalCount), nil
	}
}

func (manga *Manga) IsFavorited(ownerId primitive.ObjectID) (bool, error) {
	err := manga.GetItemFromObjectId(manga.Id)
	if err != nil {
		return false, err
	}

	for _, owner := range manga.FollowedUsers {
		if owner == ownerId {
			return true, nil
		}
	}
	return false, nil
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
