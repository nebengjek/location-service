package handlers

import (
	"location-service/bin/middlewares"
	driver "location-service/bin/modules/driver"
	"location-service/bin/modules/driver/models"
	"location-service/bin/pkg/utils"

	"github.com/labstack/echo/v4"
)

type driverHttpHandler struct {
	driverUsecaseQuery   driver.UsecaseQuery
	driverUseCaseCommand driver.UsecaseCommand
}

func InitDriverHttpHandler(e *echo.Echo, uq driver.UsecaseQuery, uc driver.UsecaseCommand) {

	handler := &driverHttpHandler{
		driverUsecaseQuery:   uq,
		driverUseCaseCommand: uc,
	}
	route := e.Group("/driver")
	route.POST("/v1/activate-beacon", handler.ActivateBeacon, middlewares.VerifyBearer)

}

func (u driverHttpHandler) ActivateBeacon(c echo.Context) error {
	var request models.BeaconRequest
	if err := c.Bind(&request); err != nil {
		return utils.ResponseError(err, c)
	}

	if err := request.Validate(); err != nil {
		return utils.ResponseError(err, c)
	}

	userId := utils.ConvertString(c.Get("userId"))
	result := u.driverUseCaseCommand.ActivateBeacon(userId, request, c.Request().Context())

	if result.Error != nil {
		return utils.ResponseError(result.Error, c)
	}

	return utils.Response(result.Data, "update beacon", 200, c)
}
