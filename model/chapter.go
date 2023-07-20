package model

import (
	"context"
	"magna/database"
	"magna/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Chapter struct {
	Id         primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Manga      primitive.ObjectID   `bson:"manga" json:"manga"`
	Name       string               `bson:"name" json:"title"`
	Cover      string               `bson:"cover" json:"cover"`
	Price      uint                 `bson:"price" json:"price"`
	UpdateTime uint                 `bson:"updateTime" json:"updateTime"`
	Images     []string             `bson:"images" json:"images"`
	OwnedUsers []primitive.ObjectID `bson:"ownedUsers" json:"ownedUsers"`
	Reports    []primitive.ObjectID `bson:"reports" json:"reports"`
}

func (chapter *Chapter) InsertToDatabase() (primitive.ObjectID, error) {
	chapter.UpdateTime = uint(time.Now().Unix())

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

func (chapter *Chapter) GetItemList(position, count int, sortField, sortType string) ([]Chapter, int, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Chapter, 0)
	filter := bson.D{{}}
	opts := options.Find()
	opts.SetSkip(int64(position))
	if sortField == "id" {
		sortField = "_id"
	}
	if sortType == "ASC" {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): 1})
	} else {
		opts.SetSort(bson.M{utils.FirstLetterToLower(sortField): -1})
	}
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

func (chapter *Chapter) GetItemFromObjectId(objID primitive.ObjectID) error {
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

func (chapter *Chapter) GetMangaChapterList(objID primitive.ObjectID) ([]Chapter, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Chapter, 0)
	multiFilter := bson.M{"manga": objID}
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

func (chapter *Chapter) GetMangaNewestChapterList(objID primitive.ObjectID, count int) ([]Chapter, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Chapter, 0)
	multiFilter := bson.M{"manga": objID}
	multiFindOpts := options.Find().SetSort(bson.D{{"updateTime", -1}})
	if count != 0 {
		multiFindOpts.SetLimit(int64(count))
	}
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

func (chapter *Chapter) IsOwned(ownerId primitive.ObjectID) (bool, error) {
	err := chapter.GetItemFromObjectId(chapter.Id)
	if err != nil {
		return false, err
	}

	if chapter.Price == 0 {
		return true, err
	}
	for _, owner := range chapter.OwnedUsers {
		if owner == ownerId {
			return true, nil
		}
	}
	return false, nil
}

func (chapter *Chapter) GetItemListFromObjectIdGroupByManga(objID []primitive.ObjectID) ([]Chapter, error) {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return nil, err
	}

	listItem := make([]Chapter, 0)
	aggregatePipeline := bson.A{}
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$match",
				bson.D{
					{"_id",
						bson.D{
							{"$in",
								objID,
							},
						},
					},
				},
			},
		},
	)
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "Manga"},
					{"localField", "manga"},
					{"foreignField", "_id"},
					{"as", "output"},
				},
			},
		},
	)
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$addFields",
				bson.D{
					{"order",
						bson.D{
							{"$indexOfArray",
								bson.A{
									bson.D{{"$first", "$output.chapters"}},
									"$_id",
								},
							},
						},
					},
				},
			},
		},
	)
	aggregatePipeline = append(aggregatePipeline,
		bson.D{
			{"$sort",
				bson.D{
					{"output.0.name", 1},
					{"order", 1},
				},
			},
		},
	)
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

func (chapter *Chapter) AddReport(reportId primitive.ObjectID) error {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	chapter.Reports = append(chapter.Reports, reportId)
	filter := bson.D{{"_id", chapter.Id}}
	update := bson.D{{"$set", bson.D{{"reports", chapter.Reports}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (chapter *Chapter) AddOwnedUsers(userId primitive.ObjectID) error {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	chapter.OwnedUsers = append(chapter.OwnedUsers, userId)
	filter := bson.D{{"_id", chapter.Id}}
	update := bson.D{{"$set", bson.D{{"ownedUsers", chapter.OwnedUsers}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (chapter *Chapter) RemoveReport(reportId primitive.ObjectID) error {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	chapter.Reports = utils.RemoveElementFromObjectIdArray(chapter.Reports, reportId)
	filter := bson.D{{"_id", chapter.Id}}
	update := bson.D{{"$set", bson.D{{"reports", chapter.Reports}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (chapter *Chapter) RemoveOwnedUsers(userId primitive.ObjectID) error {
	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	chapter.OwnedUsers = utils.RemoveElementFromObjectIdArray(chapter.OwnedUsers, userId)
	filter := bson.D{{"_id", chapter.Id}}
	update := bson.D{{"$set", bson.D{{"ownedUsers", chapter.OwnedUsers}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (chapter *Chapter) DeleteChapterById(id primitive.ObjectID) error {
	err := chapter.GetItemFromObjectId(id)
	if err != nil {
		return err
	}

	coll, err := database.GetChapterCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", chapter.Id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	for _, userId := range chapter.OwnedUsers {
		user := new(User)
		err = user.GetItemFromObjectId(userId)
		if err == nil {
			user.RemoveOwnedChapterById(chapter.Id)
		}
	}

	for _, reportId := range chapter.Reports {
		report := new(Report)
		report.DeleteReportById(reportId)
	}

	return nil
}

func (chapter *Chapter) GetTotalCount() (int, error) {
	coll, err := database.GetChapterCollection()
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
