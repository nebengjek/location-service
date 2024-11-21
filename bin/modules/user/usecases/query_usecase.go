package usecases

import (
	"context"
	"encoding/json"
	"fmt"

	"location-service/bin/modules/user"
	"location-service/bin/modules/user/models"
	httpError "location-service/bin/pkg/http-error"
	kafkaPkgConfluent "location-service/bin/pkg/kafka/confluent"

	"location-service/bin/pkg/log"
	"location-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type queryUsecase struct {
	userRepositoryQuery user.MongodbRepositoryQuery
	redisClient         redis.UniversalClient
	kafkaProducer       kafkaPkgConfluent.Producer
}

type Response struct {
	Message string      `json:"message"`
	Driver  interface{} `json:"driver"`
}

func NewQueryUsecase(mq user.MongodbRepositoryQuery, rh redis.UniversalClient, kp kafkaPkgConfluent.Producer) user.UsecaseQuery {
	return &queryUsecase{
		userRepositoryQuery: mq,
		redisClient:         rh,
		kafkaProducer:       kp,
	}
}

func (q queryUsecase) GetUser(userId string, ctx context.Context) utils.Result {
	var result utils.Result

	queryRes := <-q.userRepositoryQuery.FindOne(userId, ctx)

	if queryRes.Error != nil {
		errObj := httpError.InternalServerError("Internal server error")
		result.Error = errObj
		return result
	}
	user := queryRes.Data.(models.User)
	result.Data = &user
	return result
}

func (q *queryUsecase) FindDriver(userId string, ctx context.Context) utils.Result {
	var result utils.Result
	key := fmt.Sprintf("USER:ROUTE:%s", userId)
	var tripPlan models.RouteSummary
	redisData, errRedis := q.redisClient.Get(ctx, key).Result()
	if errRedis != nil || redisData == "" {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("Error get data from redis: %v", errRedis)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "FindDriver", utils.ConvertString(errRedis))
		return result
	}
	err := json.Unmarshal([]byte(redisData), &tripPlan)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error unmarshal tripdata: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
		return result
	}
	walletCheck := <-q.userRepositoryQuery.Findwallet(ctx, userId)
	if walletCheck.Error != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Wallet not found: %v, Please create wallet first", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
		return result
	}
	wallet := walletCheck.Data.(models.Wallet)
	if wallet.Balance <= tripPlan.MaxPrice {
		errObj := httpError.NewBadRequest()
		errObj.Message = "insufficient balance, please topup"
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "FindDriver", "")
		return result
	}
	radius := 3.0
	drivers, err := q.redisClient.GeoRadius(ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithDist:  true,
		WithCoord: true,
		Sort:      "ASC",
	}).Result()

	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error searching drivers: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "FindDriver", utils.ConvertString(err))
		return result
	}
	posibleDriver := "No driver available. Don't worry, please try again later."
	if len(drivers) > 0 {
		kafkaData := models.RequestRide{
			UserId:       userId,
			RouteSummary: tripPlan,
		}
		marshaledData, _ := json.Marshal(kafkaData)
		log.GetLogger().Info("command_usecase", "marshaled", "kafkaProducer", utils.ConvertString(marshaledData))
		q.kafkaProducer.Publish("request-ride", marshaledData)
		posibleDriver = fmt.Sprintf("Please sit back, there are %d drivers available, we will let you know", len(drivers))
	}
	result.Data = Response{
		Message: posibleDriver,
		Driver:  drivers,
	}

	return result
}
