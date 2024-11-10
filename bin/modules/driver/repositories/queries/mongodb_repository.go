package queries

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	driver "location-service/bin/modules/driver"
	"location-service/bin/modules/driver/models"
	"location-service/bin/pkg/databases/mongodb"
	"location-service/bin/pkg/utils"
)

type queryMongodbRepository struct {
	mongoDb mongodb.MongoDBLogger
}

func NewQueryMongodbRepository(mongodb mongodb.MongoDBLogger) driver.MongodbRepositoryQuery {
	return &queryMongodbRepository{
		mongoDb: mongodb,
	}
}

func (q queryMongodbRepository) FindDriver(userId string, ctx context.Context) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var driver models.User
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &driver,
			CollectionName: "user",
			Filter: bson.M{
				"userId": userId,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}

		output <- utils.Result{
			Data: driver,
		}

	}()

	return output
}

func (q queryMongodbRepository) FindWorkLog(driverId string, date string, ctx context.Context) <-chan utils.Result {
	output := make(chan utils.Result)

	go func() {
		defer close(output)

		var workLog models.WorkLog
		err := q.mongoDb.FindOne(mongodb.FindOne{
			Result:         &workLog,
			CollectionName: "work-log",
			Filter: bson.M{
				"driverId": driverId,
				"workdate": date,
			},
		}, ctx)
		if err != nil {
			output <- utils.Result{
				Error: err,
			}
		}
		if workLog.DriverID == "" {
			output <- utils.Result{
				Error: "notfound",
			}
		}
		output <- utils.Result{
			Data: workLog,
		}

	}()

	return output
}
