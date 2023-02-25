package tests

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestExtractUserIdFromEmptySessionKey(t *testing.T) {
	t.Log(primitive.ObjectIDFromHex(""))
}
