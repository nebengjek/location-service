package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httpError "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/log"

	"github.com/labstack/echo/v4"
)

// Result common output
type Result struct {
	Data  interface{}
	Error interface{}
}

// BaseWrapperModel data structure
type BaseWrapperModel struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Meta    interface{} `json:"meta,omitempty"`
}

type Meta struct {
	Method        string    `json:"method"`
	Url           string    `json:"url"`
	Code          string    `json:"code"`
	ContentLength int64     `json:"content_length"`
	Date          time.Time `json:"date"`
	Ip            string    `json:"ip"`
}

// Response function
func Response(data interface{}, message string, code int, c echo.Context) error {
	success := false
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Request().Method,
		Code:          fmt.Sprintf("%v", http.StatusOK),
		ContentLength: c.Request().ContentLength,
		Ip:            c.RealIP(),
	}
	byteMeta, _ := json.Marshal(meta)
	log.GetLogger().Info("service-info", "Logging service...", "audit-log", string(byteMeta))

	if code < http.StatusBadRequest {
		success = true
	}

	result := BaseWrapperModel{
		Success: success,
		Data:    data,
		Message: message,
		Code:    code,
	}

	return c.JSON(code, result)
}

// ResponseError function
func ResponseError(err interface{}, c echo.Context) error {
	errObj := getErrorStatusCode(err)
	result := BaseWrapperModel{
		Success: false,
		Data:    errObj.Data,
		Message: errObj.Message,
		Code:    errObj.Code,
	}
	meta := Meta{
		Date:          time.Now(),
		Url:           c.Path(),
		Method:        c.Request().Method,
		Code:          fmt.Sprintf("%v", getErrorStatusCode(err)),
		Ip:            c.RealIP(),
		ContentLength: c.Request().ContentLength,
	}
	byteMeta, _ := json.Marshal(meta)

	log.GetLogger().Error("service-error", "Logging service...", "audit-log", string(byteMeta))

	return c.JSON(errObj.ResponseCode, result)
}

func getErrorStatusCode(err interface{}) httpError.CommonErrorData {
	errData := httpError.CommonErrorData{}

	switch obj := err.(type) {
	case httpError.BadRequestData:
		errData.ResponseCode = http.StatusBadRequest
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.UnauthorizedData:
		errData.ResponseCode = http.StatusUnauthorized
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.ConflictData:
		errData.ResponseCode = http.StatusConflict
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.NotFoundData:
		errData.ResponseCode = http.StatusNotFound
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	case httpError.InternalServerErrorData:
		errData.ResponseCode = http.StatusInternalServerError
		errData.Code = obj.Code
		errData.Data = obj.Data
		errData.Message = obj.Message
		return errData
	default:
		errData.Code = http.StatusConflict
		return errData
	}
}
