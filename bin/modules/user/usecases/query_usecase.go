package usecases

import (
	"context"

	"location-service/bin/modules/user"
	"location-service/bin/modules/user/models"
	errors "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/utils"

	"github.com/redis/go-redis/v9"
)

type queryUsecase struct {
	userRepositoryQuery user.MongodbRepositoryQuery
	redisClient         redis.UniversalClient
}

func NewQueryUsecase(mq user.MongodbRepositoryQuery, rh redis.UniversalClient) user.UsecaseQuery {
	return &queryUsecase{
		userRepositoryQuery: mq,
		redisClient:         rh,
	}
}

func (q queryUsecase) GetUser(userId string, ctx context.Context) utils.Result {
	var result utils.Result

	queryRes := <-q.userRepositoryQuery.FindOne(userId, ctx)

	if queryRes.Error != nil {
		errObj := errors.InternalServerError("Internal server error")
		result.Error = errObj
		return result
	}
	user := queryRes.Data.(models.User)
	result.Data = &user
	return result
}
