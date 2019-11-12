package database

import (
	"context"

	"github.com/singkorn/jartown-services-common/neomongo"
	"github.com/singkorn/jartown-services-main/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ListReviewDatabase interface {
	GetListReview(id string) (model.ListReview, error)
}

type listReviewDatabaseImpl struct {
	listReviewsCollection neomongo.Collection
}

func GetListReviewDatabase() ListReviewDatabase {
	listReviewsCollection := database.Collection("listingsAndReviews")

	db := &listReviewDatabaseImpl{
		listReviewsCollection: neomongo.MakeCollectionImpl(listReviewsCollection),
	}

	return db
}

func (db *listReviewDatabaseImpl) GetListReview(id string) (model.ListReview, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.ListReview{}, returnError(err)
	}

	var listReview model.ListReview
	res := db.listReviewsCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	err = res.Decode(&listReview)
	if err != nil {
		return model.ListReview{}, returnError(err)
	}

	listReview.ID = id

	return listReview, nil
}
