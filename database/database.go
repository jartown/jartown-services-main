package database

import (
	"context"
	"log"

	common "github.com/singkorn/jartown-services-common"
	"github.com/singkorn/jartown-services-main/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func Init() {
	conf := config.Conf
	clientOptions := options.Client()
	clientOptions.ApplyURI(conf.MongoDB.URI)
	if conf.MongoDB.Username != "" {
		clientOptions.SetAuth(
			options.Credential{
				Username:   conf.MongoDB.Username,
				Password:   conf.MongoDB.Password,
				AuthSource: conf.MongoDB.Schema,
			},
		)
	}

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	database = client.Database(conf.MongoDB.Schema)
}

func returnError(err error) error {
	if err == mongo.ErrNoDocuments {
		return common.ErrItemNotFound
	}
	if err == primitive.ErrInvalidHex {
		return common.ErrInputValidationFailed
	}
	return err
}
