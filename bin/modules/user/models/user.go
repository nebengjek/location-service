package models

import "github.com/go-playground/validator/v10"

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
}

type LocationSuggestionRequest struct {
	CurrentLocation LocationRequest `json:"currentLocation" `
	Destination     LocationRequest `json:"destination"`
}

type Route struct {
	Origin      LocationRequest `json:"origin" `
	Destination LocationRequest `json:"destination"`
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
}
