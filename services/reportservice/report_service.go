package reportservice

import (
	"magna/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateNewReport(userId, chapterId, reportContent string) error {
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}
	chapterObjId, err := primitive.ObjectIDFromHex(chapterId)
	if err != nil {
		return err
	}
	report := new(model.Report)
	report.Chapter = chapterObjId
	report.User = userObjId
	report.Content = reportContent

	_, err = report.CreateNewReport()

	return err
}

func GetUserReport(userId string) ([]model.Report, error) {
	userObjId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	return new(model.Report).GetUserReport(userObjId)
}

func GetTotalCount() (int, error) {
	return new(model.Report).GetTotalCount()
}
