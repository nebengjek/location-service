package commands

import (
	"context"

	user "location-service/bin/modules/driver"
	"location-service/bin/modules/driver/models"
	"location-service/bin/pkg/databases/mongodb"
	"location-service/bin/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type commandMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewCommandMongodbRepository(mongodb mongodb.MongoDBLogger) user.MongodbRepositoryCommand {
	return &commandMongodbRepository{
		mongoDb: mongodb,
	}
}

func (c commandMongodbRepository) NewObjectID(ctx context.Context) string {
	return primitive.NewObjectID().Hex()
}

func (c commandMongodbRepository) UpsertBeacon(data models.WorkLog, ctx context.Context) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)
		err := c.mongoDb.UpsertOne(mongodb.UpsertOne{
			CollectionName: "work-log",
			Filter: bson.M{
				"driverId": data.DriverID,
				"workdate": data.WorkDate,
			},
			Document: bson.M{
				"driverId": data.DriverID,
				"workdate": data.WorkDate,
				"log":      data.Log,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: nil,
		}

	}()

	return output
}
