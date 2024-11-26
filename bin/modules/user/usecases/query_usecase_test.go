package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"location-service/bin/modules/user/models"
	httpError "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/utils"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMongodbRepositoryQuery struct {
	mock.Mock
}

type MockRedisClient struct {
	mock.Mock
	redis.UniversalClient
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockMongodbRepositoryQuery) FindOne(id string, ctx context.Context) <-chan utils.Result {
	args := m.Called(id, ctx)
	resultChan := make(chan utils.Result, 1)
	resultChan <- args.Get(0).(utils.Result)
	return resultChan
}

func (m *MockMongodbRepositoryQuery) Findwallet(ctx context.Context, id string) <-chan utils.Result {
	args := m.Called(ctx, id)
	resultChan := make(chan utils.Result, 1)
	resultChan <- args.Get(0).(utils.Result)
	return resultChan
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	return args.Get(0).(*redis.StringCmd)
}

func (m *MockRedisClient) GeoRadius(ctx context.Context, key string, longitude, latitude float64, query *redis.GeoRadiusQuery) *redis.GeoLocationCmd {
	args := m.Called(ctx, key, longitude, latitude, query)
	return args.Get(0).(*redis.GeoLocationCmd)
}

func (m *MockKafkaProducer) Publish(topic string, message []byte) {
	m.Called(topic, message)
}

// GetUser tests
func TestGetUser_Success(t *testing.T) {
	mockQuery := new(MockMongodbRepositoryQuery)
	mockRedis := new(MockRedisClient)
	mockKafka := new(MockKafkaProducer)

	usecase := NewQueryUsecase(mockQuery, mockRedis, mockKafka)

	ctx := context.Background()
	userId := "user123"
	expectedUser := models.User{
		Id:       userId,
		FullName: "John Doe",
		Email:    "john@example.com",
	}

	mockQuery.On("FindOne", userId, ctx).Return(utils.Result{Data: expectedUser})

	result := usecase.GetUser(userId, ctx)

	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Data)
	user := result.Data.(*models.User)
	assert.Equal(t, expectedUser.Id, user.Id)
	assert.Equal(t, expectedUser.FullName, user.FullName)
	assert.Equal(t, expectedUser.Email, user.Email)
}

func TestGetUser_NotFound(t *testing.T) {
	mockQuery := new(MockMongodbRepositoryQuery)
	mockRedis := new(MockRedisClient)
	mockKafka := new(MockKafkaProducer)

	usecase := NewQueryUsecase(mockQuery, mockRedis, mockKafka)

	ctx := context.Background()
	userId := "nonexistent"

	mockQuery.On("FindOne", userId, ctx).Return(utils.Result{Error: errors.New("user not found")})

	result := usecase.GetUser(userId, ctx)

	assert.NotNil(t, result.Error)
	if err, ok := result.Error.(httpError.NotFoundData); ok {
		assert.Equal(t, httpError.NewInternalServerError().Code, err.Code)
	} else {
		t.Errorf("expected error of type httpError.HttpError, got %T", result.Error)
	}
}

// FindDriver tests
func TestFindDriver_Success(t *testing.T) {
	mockQuery := new(MockMongodbRepositoryQuery)
	mockRedis := new(MockRedisClient)
	mockKafka := new(MockKafkaProducer)

	usecase := NewQueryUsecase(mockQuery, mockRedis, mockKafka)

	ctx := context.Background()
	userId := "user123"
	key := "USER:ROUTE:user123"
	tripPlan := models.RouteSummary{
		MaxPrice: 1000,
		Route: models.Route{
			Origin: models.LocationRequest{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		},
	}
	tripPlanData, _ := json.Marshal(tripPlan)
	wallet := models.Wallet{
		Balance: 2000,
	}

	drivers := []redis.GeoLocation{{Name: "driver1"}}
	mockRedis.On("Get", ctx, key).Return(redis.NewStringResult(string(tripPlanData), nil))
	mockQuery.On("Findwallet", ctx, userId).Return(utils.Result{Data: wallet})
	mockRedis.On("GeoRadius", ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, mock.Anything).Return(redis.NewGeoLocationCmdResult(drivers, nil))
	mockKafka.On("Publish", "request-ride", mock.Anything).Return(nil)

	result := usecase.FindDriver(userId, ctx)

	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Data)
	response := result.Data.(Response)
	assert.Equal(t, "Please sit back, there are 1 drivers available, we will let you know", response.Message)
	assert.Len(t, response.Driver.([]redis.GeoLocation), 1)
}

func TestFindDriver_GeoRadiusError(t *testing.T) {
	mockQuery := new(MockMongodbRepositoryQuery)
	mockRedis := new(MockRedisClient)
	mockKafka := new(MockKafkaProducer)

	usecase := NewQueryUsecase(mockQuery, mockRedis, mockKafka)

	ctx := context.Background()
	userId := "user123"
	key := "USER:ROUTE:user123"
	tripPlan := models.RouteSummary{
		MaxPrice: 1000,
		Route: models.Route{
			Origin: models.LocationRequest{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		},
	}
	tripPlanData, _ := json.Marshal(tripPlan)
	wallet := models.Wallet{
		Balance: 2000,
	}

	mockRedis.On("Get", ctx, key).Return(redis.NewStringResult(string(tripPlanData), nil))
	mockQuery.On("Findwallet", ctx, userId).Return(utils.Result{Data: wallet})
	mockRedis.On("GeoRadius", ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, mock.Anything).Return(redis.NewGeoLocationCmdResult(nil, errors.New("geo radius error")))

	result := usecase.FindDriver(userId, ctx)

	assert.NotNil(t, result.Error)
	if err, ok := result.Error.(httpError.InternalServerErrorData); ok {
		assert.Equal(t, httpError.NewInternalServerError().Code, err.Code)
	} else {
		t.Errorf("expected error of type httpError.HttpError, got %T", result.Error)
	}
}

func TestFindDriver_KafkaPublishError(t *testing.T) {
	mockQuery := new(MockMongodbRepositoryQuery)
	mockRedis := new(MockRedisClient)
	mockKafka := new(MockKafkaProducer)

	usecase := NewQueryUsecase(mockQuery, mockRedis, mockKafka)

	ctx := context.Background()
	userId := "user123"
	key := "USER:ROUTE:user123"
	tripPlan := models.RouteSummary{
		MaxPrice: 1000,
		Route: models.Route{
			Origin: models.LocationRequest{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		},
	}
	tripPlanData, _ := json.Marshal(tripPlan)
	wallet := models.Wallet{
		Balance: 2000,
	}

	drivers := []redis.GeoLocation{{Name: "driver1"}}
	mockRedis.On("Get", ctx, key).Return(redis.NewStringResult(string(tripPlanData), nil))
	mockQuery.On("Findwallet", ctx, userId).Return(utils.Result{Data: wallet})
	mockRedis.On("GeoRadius", ctx, "drivers-locations", tripPlan.Route.Origin.Longitude, tripPlan.Route.Origin.Latitude, mock.Anything).Return(redis.NewGeoLocationCmdResult(drivers, nil))
	mockKafka.On("Publish", "request-ride", mock.Anything).Return(errors.New("kafka publish error"))

	result := usecase.FindDriver(userId, ctx)

	assert.Nil(t, result.Error)
	assert.NotNil(t, result.Data)
	response := result.Data.(Response)
	assert.Equal(t, "Please sit back, there are 1 drivers available, we will let you know", response.Message)
}
