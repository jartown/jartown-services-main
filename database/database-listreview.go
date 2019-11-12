package database

import (
	"context"

	"gitlab.com/diancai/diancai-api/model"
	"gitlab.com/diancai/diancai-services-common/neomongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopDatabase interface {
	GetShop(id string) (model.Shop, error)
}

type shopDatabaseImpl struct {
	shopsCollection      neomongo.Collection
	categoriesCollection neomongo.Collection
	menusCollection      neomongo.Collection
}

func GetShopDatabase() ShopDatabase {
	shopsCollection := database.Collection("shops")
	categoriesCollection := database.Collection("categories")
	menusCollection := database.Collection("menus")

	db := &shopDatabaseImpl{
		shopsCollection:      neomongo.MakeCollectionImpl(shopsCollection),
		categoriesCollection: neomongo.MakeCollectionImpl(categoriesCollection),
		menusCollection:      neomongo.MakeCollectionImpl(menusCollection),
	}

	_, err := menusCollection.Indexes().CreateMany(context.TODO(),
		[]mongo.IndexModel{
			{Keys: bson.D{
				{Key: "shop_id", Value: 1},
			}},
		},
	)

	if err != nil {
		panic(err)
	}

	_, err = categoriesCollection.Indexes().CreateMany(context.TODO(),
		[]mongo.IndexModel{
			{Keys: bson.D{
				{Key: "shop_id", Value: 1},
				{Key: "position", Value: 1},
			}},
		},
	)

	if err != nil {
		panic(err)
	}

	return db
}

func (db *shopDatabaseImpl) GetShop(id string) (model.Shop, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return model.Shop{}, returnError(err)
	}

	var shop model.Shop
	res := db.shopsCollection.FindOne(context.TODO(), bson.D{{Key: "_id", Value: objID}})
	err = res.Decode(&shop)
	if err != nil {
		return model.Shop{}, returnError(err)
	}

	shop.ID = id

	return shop, nil
}
