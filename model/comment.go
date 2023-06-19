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

type Comment struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Manga       primitive.ObjectID `bson:"manga"`
	User        primitive.ObjectID `bson:"user"`
	Content     string             `bson:"content"`
	TimeCreated uint               `bson:"timeCreated"`
}

func (comment *Comment) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), comment)
	if err != nil {
		return [12]byte{}, err
	}

	comment.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (comment *Comment) GetItemList(position, count int, sortField, sortType string) ([]Comment, int, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]Comment, 0)
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

func (comment *Comment) GetItemFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(comment)
	if err != nil {
		return err
	}

	return nil
}

func (comment *Comment) GetListCommentFromMangaId(mangaId primitive.ObjectID, position, count int) ([]Comment, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return nil, err
	}

	filter := bson.D{primitive.E{Key: "manga", Value: mangaId}}
	opts := options.Find().SetSort(bson.D{{"timeCreated", -1}}).SetSkip(int64(position)).SetLimit(int64(count))
	cursor, err := coll.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	listItem := make([]Comment, 0)
	err = cursor.All(context.TODO(), &listItem)
	if err != nil {
		return nil, err
	}

	return listItem, err
}

func (comment *Comment) GetMangaCommentCount(mangaId primitive.ObjectID) (int, error) {
	coll, err := database.GetCommentCollection()
	if err != nil {
		return 0, err
	}

	filter := bson.D{primitive.E{Key: "manga", Value: mangaId}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (comment *Comment) CreateNewComment() (primitive.ObjectID, error) {
	comment.TimeCreated = uint(time.Now().Unix())

	coll, err := database.GetCommentCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), comment)
	if err != nil {
		return [12]byte{}, err
	}

	comment.Id = result.InsertedID.(primitive.ObjectID)

	user := new(User)
	err = user.GetItemFromObjectId(comment.User)
	if err != nil {
		return [12]byte{}, err
	}
	err = user.AddComment(comment.Id)
	if err != nil {
		return [12]byte{}, err
	}

	manga := new(Manga)
	err = manga.GetItemFromObjectId(comment.Manga)
	if err != nil {
		return [12]byte{}, err
	}
	err = manga.AddComment(comment.Id)
	if err != nil {
		return [12]byte{}, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (comment *Comment) DeleteCommentById(id primitive.ObjectID) error {
	err := comment.GetItemFromObjectId(id)
	if err != nil {
		return err
	}

	coll, err := database.GetCommentCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", comment.Id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}

	user := new(User)
	err = user.GetItemFromObjectId(comment.User)
	// If successful get user then remove comment from user
	if err == nil {
		user.RemoveComment(comment.Id)
	}

	manga := new(Manga)
	err = manga.GetItemFromObjectId(comment.Manga)
	// If successful get manga then remove comment from manga
	if err == nil {
		manga.RemoveComment(comment.Id)
	}

	return nil
}
