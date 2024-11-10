package main

import (
	"context"
	"fmt"

	"location-service/bin/pkg/log"
	"location-service/bin/pkg/redis"
	"net/http"
	"os"
	"os/signal"
	"time"

	"location-service/bin/config"
	userHandler "location-service/bin/modules/user/handlers"
	userRepoCommands "location-service/bin/modules/user/repositories/commands"
	userRepoQueries "location-service/bin/modules/user/repositories/queries"
	userUsecase "location-service/bin/modules/user/usecases"

	driverHandler "location-service/bin/modules/driver/handlers"
	driverRepoCommands "location-service/bin/modules/driver/repositories/commands"
	driverRepoQueries "location-service/bin/modules/driver/repositories/queries"
	driverUsecase "location-service/bin/modules/driver/usecases"

	"location-service/bin/pkg/apm"
	"location-service/bin/pkg/databases/mongodb"
	kafkaConfluent "location-service/bin/pkg/kafka/confluent"
	"location-service/bin/pkg/utils"

	"location-service/bin/pkg/validator"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"go.elastic.co/apm/module/apmechov4"
)

func main() {
	apm.InitConnection()
	redis.LoadConfig()
	redis.InitConnection()
	mongodb.InitConnection()
	kafkaConfluent.InitKafkaConfig()
	log.Init()
	e := echo.New()
	e.Validator = &validator.CustomValidator{Validator: validator.New()}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper:          middleware.DefaultSkipper,
		Format:           `[ROUTE] ${time_rfc3339} | ${status} | ${latency_human} ${latency} | ${method} | ${uri}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}))
	e.Use(middleware.Recover())
	e.Use(apmechov4.Middleware(apmechov4.WithTracer(apm.GetTracer())))

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))
	setConfluentEvents()

	setHttp(e)

	listenerPort := fmt.Sprintf(":%s", config.GetConfig().AppPort)
	e.Logger.Fatal(e.Start(listenerPort))

	server := &http.Server{
		Addr:         listenerPort,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  5 * time.Second,
	}
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		log.GetLogger().Info("main", "Server message-service is shutting down...", "gracefull", "")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			log.GetLogger().Info("main", fmt.Sprintf("Could not gracefully shutdown the server order-service: %v\n", err), "gracefull", "")
		}
		close(done)
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.GetLogger().Info("main", fmt.Sprintf("Could not listen on %s: %v\n", config.GetConfig().AppPort, err), "gracefull", "")
	}

	<-done
	log.GetLogger().Info("main", fmt.Sprintf("Server %s stopped", config.GetConfig().AppName), "gracefull", "")
}

func setConfluentEvents() {

}

func setHttp(e *echo.Echo) {
	redisClient := redis.GetClient()
	e.GET("/v1/health-check", func(c echo.Context) error {
		log.GetLogger().Info("main", "This service is running properly", "setConfluentEvents", "")
		return utils.Response(nil, "This service is running properly", 200, c)
	})
	kafkaProducer, err := kafkaConfluent.NewProducer(kafkaConfluent.GetConfig().GetKafkaConfig(), log.GetLogger())
	if err != nil {
		panic(err)
	}

	userQueryMongodbRepo := userRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	userCommandMongodbRepo := userRepoCommands.NewCommandMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), log.GetLogger()))

	userQueryUsecase := userUsecase.NewQueryUsecase(userQueryMongodbRepo, redisClient, kafkaProducer)
	userCommandUsecase := userUsecase.NewCommandUsecase(userQueryMongodbRepo, userCommandMongodbRepo, config.GetConfig().GoogleApiKey, redisClient)

	driverQueryMongodbRepo := driverRepoQueries.NewQueryMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetSlaveConn(), mongodb.GetSlaveDBName(), log.GetLogger()))
	driverCommandMongodbRepo := driverRepoCommands.NewCommandMongodbRepository(mongodb.NewMongoDBLogger(mongodb.GetMasterConn(), mongodb.GetMasterDBName(), log.GetLogger()))

	driverQueryUsecase := driverUsecase.NewQueryUsecase(driverQueryMongodbRepo, redisClient)
	driverCommandUsecase := driverUsecase.NewCommandUsecase(driverQueryMongodbRepo, driverCommandMongodbRepo, redisClient)

	userHandler.InituserHttpHandler(e, userQueryUsecase, userCommandUsecase)
	driverHandler.InitDriverHttpHandler(e, driverQueryUsecase, driverCommandUsecase)
}
