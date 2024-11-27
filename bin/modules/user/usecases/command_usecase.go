package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"location-service/bin/modules/user"
	"location-service/bin/modules/user/models"
	httpError "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/log"
	"location-service/bin/pkg/utils"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"googlemaps.github.io/maps"
)

type commandUsecase struct {
	userRepositoryQuery   user.MongodbRepositoryQuery
	userRepositoryCommand user.MongodbRepositoryCommand
	googleMapsAPIKey      string
	redisClient           redis.UniversalClient
}

func NewCommandUsecase(mq user.MongodbRepositoryQuery, mc user.MongodbRepositoryCommand, gkey string, rc redis.UniversalClient) user.UsecaseCommand {
	return &commandUsecase{
		userRepositoryQuery:   mq,
		userRepositoryCommand: mc,
		googleMapsAPIKey:      gkey,
		redisClient:           rc,
	}
}

func (c *commandUsecase) PostLocation(userId string, payload models.LocationSuggestionRequest, ctx context.Context) utils.Result {
	var result utils.Result
	mapsClient, err := maps.NewClient(maps.WithAPIKey(c.googleMapsAPIKey))
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("error creating Google Maps client: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PostLocation", utils.ConvertString(err))
		return result
	}

	routeSuggestion, err := c.getRouteSuggestions(ctx, mapsClient, payload.CurrentLocation, payload.Destination)
	if err != nil {
		errObj := httpError.NewNotFound()
		errObj.Message = fmt.Sprintf("error getRouteSuggestions: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PostLocation", utils.ConvertString(err))
		return result
	}
	key := fmt.Sprintf("USER:ROUTE:%s", userId)
	routeSuggestion.Route.Origin = payload.CurrentLocation
	routeSuggestion.Route.Destination = payload.Destination
	routeSummaryJSON, err := json.Marshal(routeSuggestion)
	if err != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error marshalling RouteSummary: %v", err)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PostLocation", utils.ConvertString(err))
		return result
	}

	redisErr := c.redisClient.Set(ctx, key, routeSummaryJSON, 60*time.Minute).Err()
	if redisErr != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Error saving to redis: %v", redisErr)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "PostLocation", utils.ConvertString(redisErr))
		return result
	}
	result.Data = routeSuggestion
	return result
}

func (c *commandUsecase) getRouteSuggestions(ctx context.Context, mapsClient *maps.Client, currentRequest models.LocationRequest, destinationRequest models.LocationRequest) (*models.RouteSummary, error) {
	origin := fmt.Sprintf("%f,%f", currentRequest.Latitude, currentRequest.Longitude)
	destination := fmt.Sprintf("%f,%f", destinationRequest.Latitude, destinationRequest.Longitude)
	departureTime := time.Now().Add(5 * time.Minute).Unix()

	req := &maps.DirectionsRequest{
		Origin:        origin,
		Destination:   destination,
		Mode:          maps.TravelModeDriving,
		Alternatives:  true,
		Optimize:      true,
		DepartureTime: fmt.Sprintf("%d", departureTime),
		TrafficModel:  maps.TrafficModelBestGuess,
	}

	routes, _, err := mapsClient.Directions(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error making directions request: %w", err)
	}

	if len(routes) == 0 {
		return nil, fmt.Errorf("no routes found")
	}

	const pricePerKm = 3000.0
	var minPrice, maxPrice float64
	var bestRouteKm, bestRoutePrice, bestRouteDuration float64

	minPrice = math.MaxFloat64
	maxPrice = -math.MaxFloat64

	for _, route := range routes {
		totalDistance := 0.0
		totalDuration := 0.0

		for _, leg := range route.Legs {
			totalDistance += float64(leg.Distance.Meters)
			totalDuration += float64(leg.DurationInTraffic.Seconds())
		}

		distanceInKm := totalDistance / 1000.0
		price := distanceInKm * pricePerKm

		if price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}

		if bestRouteKm == 0 || price < bestRoutePrice {
			bestRouteKm = distanceInKm
			bestRoutePrice = price
			bestRouteDuration = totalDuration / 60
		}
	}

	return &models.RouteSummary{
		MinPrice:          minPrice,
		MaxPrice:          maxPrice,
		BestRouteKm:       bestRouteKm,
		BestRoutePrice:    bestRoutePrice,
		BestRouteDuration: utils.FormatDuration(int(math.Ceil(bestRouteDuration))),
		Duration:          int(math.Ceil(bestRouteDuration)),
	}, nil

}
