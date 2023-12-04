package model

import (
	"context"
	"errors"
	"magna/database"
	"magna/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id            primitive.ObjectID         `bson:"_id,omitempty" json:"id"`
	Email         string                     `bson:"email" json:"email"`
	Password      string                     `bson:"password" json:"password"`
	Role          string                     `bson:"role" json:"role"`
	DisplayName   string                     `bson:"displayName" json:"displayname"`
	Avatar        string                     `bson:"avatar" json:"avatar"`
	FirstName     string                     `bson:"firstName" json:"firstname"`
	LastName      string                     `bson:"lastName" json:"lastname"`
	Gender        int                        `bson:"gender" json:"gender"`
	FollowMangas  []primitive.ObjectID       `bson:"followMangas" json:"followMangas"`
	OwnedChapters []primitive.ObjectID       `bson:"ownedChapters" json:"ownedChapters"`
	Comments      []primitive.ObjectID       `bson:"comments" json:"comments"`
	Reports       []primitive.ObjectID       `bson:"reports" json:"reports"`
	Rate          map[primitive.ObjectID]int `bson:"rate" json:"rate"`
}

func (user *User) InsertToDatabase() (primitive.ObjectID, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return [12]byte{}, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (user *User) UpdateInfo() error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", user.Id}}
	// fmt.Println("Hello", user.Gender)
	update := bson.D{{"$set", bson.D{
		{"displayName", user.DisplayName},
		{"avatar", user.Avatar},
		{"firstName", user.FirstName},
		{"lastName", user.LastName},
		{"gender", user.Gender},
	}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) GetItemFromObjectId(objID primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "_id", Value: objID}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) GetItemList(position, count int, sortField, sortType string) ([]User, int, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]User, 0)
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

func (user *User) GetItemListFilterByDisplayName(position, count int, sortField, sortType, filterValue string) ([]User, int, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]User, 0)
	filter := bson.D{{"displayName", primitive.Regex{Pattern: filterValue, Options: "i"}}}
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

func (user *User) GetItemListFilterByEmail(position, count int, sortField, sortType, filterValue string) ([]User, int, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]User, 0)
	filter := bson.D{{"email", primitive.Regex{Pattern: filterValue, Options: "i"}}}
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

func (user *User) GetItemListFilterByFollowManga(position, count int, sortField, sortType, filterValue string) ([]User, int, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return nil, 0, err
	}

	filterValueObjId, err := primitive.ObjectIDFromHex(filterValue)
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]User, 0)
	filter := bson.M{"followMangas": filterValueObjId}
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

func (user *User) GetItemListFilterByOwnedChapters(position, count int, sortField, sortType, filterValue string) ([]User, int, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return nil, 0, err
	}

	filterValueObjId, err := primitive.ObjectIDFromHex(filterValue)
	if err != nil {
		return nil, 0, err
	}

	listItem := make([]User, 0)
	filter := bson.M{"ownedChapters": filterValueObjId}
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

func (user *User) GetItemFromEmail(email string) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	opts := options.FindOne()
	found := coll.FindOne(context.TODO(), filter, opts)
	if found.Err() != nil {
		return found.Err()
	}
	err = found.Decode(user)
	if err != nil {
		return err
	}

	return nil
}

func (user *User) CreateNewUser(email, password, role string) (primitive.ObjectID, error) {
	existed, err := checkExistedEmail(email)
	if err != nil {
		return [12]byte{}, err
	}
	if existed {
		return [12]byte{}, errors.New("Email đã tồn tại")
	}

	user.Email = email
	user.Password, _ = utils.Hash(password)
	user.Role = role
	user.Avatar = "https://st3.depositphotos.com/1767687/16607/v/450/depositphotos_166074422-stock-illustration-default-avatar-profile-icon-grey.jpg"
	user.DisplayName = email
	user.FirstName = "Tên"
	user.LastName = "Họ"
	user.Gender = 0
	user.Rate = make(map[primitive.ObjectID]int)

	coll, err := database.GetUserCollection()
	if err != nil {
		return [12]byte{}, err
	}

	result, err := coll.InsertOne(context.TODO(), user)
	if err != nil {
		return [12]byte{}, err
	}

	user.Id = result.InsertedID.(primitive.ObjectID)
	return result.InsertedID.(primitive.ObjectID), nil
}

func (user *User) AddComment(commentId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.Comments = append(user.Comments, commentId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"comments", user.Comments}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) AddReport(reportId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.Reports = append(user.Reports, reportId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"reports", user.Reports}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) DeleteUserById(id primitive.ObjectID) error {
	err := user.GetItemFromObjectId(id)
	if err != nil {
		return err
	}

	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", user.Id}}
	_, err = coll.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	for _, mangaId := range user.FollowMangas {
		manga := new(Manga)
		err = manga.GetItemFromObjectId(mangaId)
		// If successful get manga then remove folllowuser from manga
		if err == nil {
			manga.RemoveFollowUser(user.Id)
		}
	}
	for _, chapterId := range user.OwnedChapters {
		chapter := new(Chapter)
		err = chapter.GetItemFromObjectId(chapterId)
		// If successful get chapter then remove ownedUser from chapter
		if err == nil {
			chapter.RemoveOwnedUsers(user.Id)
		}
	}
	for _, commentId := range user.Comments {
		new(Comment).DeleteCommentById(commentId)
	}
	for _, reportId := range user.Reports {
		new(Report).DeleteReportById(reportId)
	}

	return nil
}

func (user *User) RemoveReport(reportId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.Reports = utils.RemoveElementFromObjectIdArray(user.Reports, reportId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"reports", user.Reports}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RemoveComment(commentId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.Comments = utils.RemoveElementFromObjectIdArray(user.Comments, commentId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"comments", user.Comments}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RemoveFollowedMangaById(mangaId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.FollowMangas = utils.RemoveElementFromObjectIdArray(user.FollowMangas, mangaId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"followMangas", user.FollowMangas}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RemoveOwnedChapterById(chapterId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	user.OwnedChapters = utils.RemoveElementFromObjectIdArray(user.OwnedChapters, chapterId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"ownedChapters", user.OwnedChapters}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SetRate() error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"rate", user.Rate}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SetOwned() error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"ownedChapters", user.OwnedChapters}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) RemoveRateManga(mangaId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	delete(user.Rate, mangaId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"rate", user.Rate}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) SetFavoriteManga(mangaId primitive.ObjectID) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}
	for index, followMangas := range user.FollowMangas {
		if followMangas == mangaId {
			ret := make([]primitive.ObjectID, 0)
			ret = append(ret, user.FollowMangas[:index]...)
			user.FollowMangas = append(ret, user.FollowMangas[index+1:]...)
			filter := bson.D{{"_id", user.Id}}
			update := bson.D{{"$set", bson.D{{"followMangas", user.FollowMangas}}}}
			_, err = coll.UpdateOne(context.TODO(), filter, update)
			if err != nil {
				return err
			}
			return nil
		}
	}

	user.FollowMangas = append(user.FollowMangas, mangaId)
	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{"followMangas", user.FollowMangas}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) GetTotalCount() (int, error) {
	coll, err := database.GetUserCollection()
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

func (user *User) IsAdmin() (bool, error) {
	err := user.GetItemFromObjectId(user.Id)
	if err != nil {
		return false, err
	}

	return user.Role == "Quản trị viên", nil
}

func (user *User) Update(fieldName string, fieldValue interface{}) error {
	coll, err := database.GetUserCollection()
	if err != nil {
		return err
	}

	filter := bson.D{{"_id", user.Id}}
	update := bson.D{{"$set", bson.D{{fieldName, fieldValue}}}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

func checkExistedEmail(email string) (bool, error) {
	coll, err := database.GetUserCollection()
	if err != nil {
		return false, err
	}

	filter := bson.D{primitive.E{Key: "email", Value: email}}
	count, err := coll.CountDocuments(context.TODO(), filter)
	if err != nil {
		return false, err
	}

	if count != 0 {
		return true, nil
	}
	return false, nil
}
