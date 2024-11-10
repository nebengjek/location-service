package apm

import (
	"location-service/bin/config"
	"os"

	"go.elastic.co/apm"
	"go.elastic.co/apm/transport"
)

func InitConnection() {
	os.Setenv("ELASTIC_APM_SERVER_URL", config.GetConfig().APMUrl)
	os.Setenv("ELASTIC_APM_SECRET_TOKEN", config.GetConfig().APMSecretToken)

	if _, err := transport.InitDefault(); err != nil {
		panic(err)
	}
}

func GetTracer() *apm.Tracer {
	tracer, _ := apm.NewTracer(config.GetConfig().AppName, config.GetConfig().AppVersion)
	return tracer
}
