package user

import (
	"context"

	"location-service/bin/modules/driver/models"
	"location-service/bin/pkg/utils"
	//"go.mongodb.org/mongo-driver/bson"
)

type UsecaseQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
}

type UsecaseCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	ActivateBeacon(userId string, payload models.BeaconRequest, ctx context.Context) utils.Result
}

type MongodbRepositoryQuery interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	FindWorkLog(driverId string, date string, ctx context.Context) <-chan utils.Result
	FindDriver(userId string, ctx context.Context) <-chan utils.Result
}

type MongodbRepositoryCommand interface {
	// idiomatic go, ctx first before payload. See https://pkg.go.dev/context#pkg-overview
	NewObjectID(ctx context.Context) string
	UpsertBeacon(data models.WorkLog, ctx context.Context) <-chan utils.Result
}
