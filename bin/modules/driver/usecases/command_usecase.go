package usecases

import (
	"context"
	"fmt"
	"location-service/bin/config"

	driver "location-service/bin/modules/driver"
	"location-service/bin/modules/driver/models"
	httpError "location-service/bin/pkg/http-error"
	"location-service/bin/pkg/log"
	"location-service/bin/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type commandUsecase struct {
	driverRepositoryQuery   driver.MongodbRepositoryQuery
	driverRepositoryCommand driver.MongodbRepositoryCommand
	redisClient             redis.UniversalClient
}

func NewCommandUsecase(mq driver.MongodbRepositoryQuery, mc driver.MongodbRepositoryCommand, rc redis.UniversalClient) driver.UsecaseCommand {
	return &commandUsecase{
		driverRepositoryQuery:   mq,
		driverRepositoryCommand: mc,
		redisClient:             rc,
	}
}

func (c *commandUsecase) ActivateBeacon(driverId string, payload models.BeaconRequest, ctx context.Context) utils.Result {
	var result utils.Result
	driverInfo := <-c.driverRepositoryQuery.FindDriver(driverId, ctx)
	if driverInfo.Error != nil {
		errObj := httpError.BadRequest("Profile Driver not completed")
		result.Error = errObj
		return result
	}
	now := time.Now()
	formattedDate := now.Format("2006-01-02")
	driver, _ := driverInfo.Data.(models.User)
	var workLogData models.WorkLog

	workLog := <-c.driverRepositoryQuery.FindWorkLog(driver.Id, formattedDate, ctx)
	if workLog.Error == nil && workLog.Data != nil {
		workLogData = workLog.Data.(models.WorkLog)
	} else {
		workLogData = models.WorkLog{
			DriverID: driver.Id,
			WorkDate: formattedDate,
		}
	}
	if len(workLogData.Log) > 0 {
		lastLog := workLogData.Log[len(workLogData.Log)-1]

		if lastLog.Status == payload.Status || now.Sub(lastLog.WorkTime) < time.Minute {
			fmt.Println("Log tidak diperbarui karena status sama")
		} else {
			workLogData.Log = append(workLogData.Log, models.LogActivity{
				WorkTime: now,
				Active:   payload.Status == "work",
				Status:   payload.Status,
			})
		}
	} else {
		workLogData.Log = append(workLogData.Log, models.LogActivity{
			WorkTime: now,
			Active:   payload.Status == "work",
			Status:   payload.Status,
		})
	}
	beacon := <-c.driverRepositoryCommand.UpsertBeacon(workLogData, ctx)
	if beacon.Error != nil {
		errObj := httpError.NewInternalServerError()
		errObj.Message = fmt.Sprintf("Failed update worklog: %v", beacon.Error)
		result.Error = errObj
		log.GetLogger().Error("command_usecase", errObj.Message, "UpsertBeacon", utils.ConvertString(beacon.Error))
		return result
	}
	urlSocket := fmt.Sprintf("%s?driver=%s", config.GetConfig().SocketUrl, driver.Id)
	if payload.Status != "work" {
		urlSocket = "selamat istirahat"
	}

	result.Data = urlSocket
	return result
}
