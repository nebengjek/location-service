package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id           string `json:"_id" bson:"_id"`
	Email        string `json:"email" bson:"email" validate:"required,email"`
	FullName     string `json:"fullName" bson:"fullName" validate:"required,min=3,max=100"`
	MobileNumber string `json:"mobileNumber" bson:"mobileNumber" validate:"required"`
}

type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Keyword   string  `json:"keyword"`
}

type LocationSuggestion struct {
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	StreetName   string  `json:"streetName"`
	NameLocation string  `json:"nameLocation"`
}

type LocationSuggestionResponse struct {
	CurrentLocation []LocationSuggestion `json:"currentLocation"`
	Destination     []LocationSuggestion `json:"destination"`
}

type LocationRequest struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Address   string  `json:"address" validate:"required"`
}

type LocationSuggestionRequest struct {
	CurrentLocation LocationRequest `json:"currentLocation" validate:"required"`
	Destination     LocationRequest `json:"destination" validate:"required"`
}

type Route struct {
	Origin      LocationRequest `json:"origin" `
	Destination LocationRequest `json:"destination"`
}

type RequestRide struct {
	RouteSummary RouteSummary `json:"routeSummary" bson:"routeSummary"`
	UserId       string       `json:"userId" bson:"userId"`
}

func (r *LocationSuggestionRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type RouteSummary struct {
	Route             Route   `json:"route"`
	MinPrice          float64 `json:"minPrice"`
	MaxPrice          float64 `json:"maxPrice"`
	BestRouteKm       float64 `json:"bestRouteKm"`
	BestRoutePrice    float64 `json:"bestRoutePrice"`
	BestRouteDuration string  `json:"bestRouteDuration"`
	Duration          int     `json:"duration"`
}

type Wallet struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID         string             `bson:"userId" json:"userId"`
	Balance        float64            `bson:"balance" json:"balance"`
	TransactionLog []TransactionLog   `bson:"transactionLog" json:"transactionLog"`
	LastUpdated    time.Time          `bson:"lastUpdated" json:"lastUpdated"`
}

type TransactionLog struct {
	TransactionID string    `bson:"transactionId" json:"transactionId"`
	Amount        float64   `bson:"amount" json:"amount"`
	Type          string    `bson:"type" json:"type"`
	Description   string    `bson:"description" json:"description"`
	Timestamp     time.Time `bson:"timestamp" json:"timestamp"`
}
