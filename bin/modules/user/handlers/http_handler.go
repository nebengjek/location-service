package handlers

import (
	"fmt"
	"location-service/bin/middlewares"
	"location-service/bin/modules/user"
	"location-service/bin/modules/user/models"
	httpError "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

type userHttpHandler struct {
	userUsecaseQuery   user.UsecaseQuery
	userUseCaseCommand user.UsecaseCommand
}

func InituserHttpHandler(e *echo.Echo, uq user.UsecaseQuery, uc user.UsecaseCommand) {

	handler := &userHttpHandler{
		userUsecaseQuery:   uq,
		userUseCaseCommand: uc,
	}
	route := e.Group("/users")
	route.GET("/profile", handler.Getuser, middlewares.VerifyBearer)
	route.POST("/v1/post-location", handler.PostLocation, middlewares.VerifyBearer)
	route.GET("/v1/find-driver", handler.FindDriver, middlewares.VerifyBearer)

}

func (u userHttpHandler) Getuser(c echo.Context) error {

	userId := utils.ConvertString(c.Get("userId"))
	result := u.userUsecaseQuery.GetUser(userId, c.Request().Context())

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Get user success", 200, c)
}

func (u userHttpHandler) PostLocation(c echo.Context) error {
	var request models.LocationSuggestionRequest
	if err := c.Bind(&request); err != nil {
		return utils.ResponseError(err, c)
	}

	if err := request.Validate(); err != nil {
		errObj := httpError.NewBadRequest()
		errObj.Message = fmt.Sprintf("Request validation error: %v", err.Error())
		return utils.ResponseError(errObj, c)
	}

	userId := utils.ConvertString(c.Get("userId"))
	result := u.userUseCaseCommand.PostLocation(userId, request, c.Request().Context())

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "Route Estimation success", 200, c)
}

func (u userHttpHandler) FindDriver(c echo.Context) error {
	userId := utils.ConvertString(c.Get("userId"))
	result := u.userUsecaseQuery.FindDriver(userId, c.Request().Context())

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "finding driver", 200, c)
}
