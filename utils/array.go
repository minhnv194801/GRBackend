package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func RemoveEmptyElementsFromStringArray(arr []string) []string {
	var r []string
	for _, str := range arr {
		if str != "" {
			r = append(r, str)
		}
	}
	arr = r
	return r
}

func RemoveElementFromObjectIdArray(arr []primitive.ObjectID, removeElement primitive.ObjectID) []primitive.ObjectID {
	var r []primitive.ObjectID
	for _, objId := range arr {
		if objId.Hex() != removeElement.Hex() {
			r = append(r, objId)
		}
	}
	arr = r
	return r
}
