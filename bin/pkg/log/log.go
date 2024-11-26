package log

import (
	"fmt"
	"runtime"

	"location-service/bin/config"
	// "location-service/bin/pkg/logstash"
)

type Log struct {
	appName  string
	logLevel int
}

var logger Log

// var logstsh *logstash.Logstash
var mapOfLogLevel = map[string]int{
	"ERROR": 2,
	"DEBUG": 1,
}
var serviceName map[string]interface{}

func Init() {
	// logstsh = logstash.New(config.GetConfig().LogstashHost, config.GetConfig().LogstashPortInt(), 5)
	logger = Log{
		appName:  config.GetConfig().AppName,
		logLevel: mapOfLogLevel[config.GetConfig().LogLevel],
	}

	// if _, err := logstsh.Connect(); err != nil {
	// 	println(err.Error())
	// }
}

func GetLogger() Log {
	return logger
}

// type logstashPayload struct {
// 	ServiceName interface{} `json:"serviceName"`
// 	Context     string      `json:"context"`
// 	Scope       string      `json:"scope"`
// 	Message     string      `json:"message"`
// 	Meta        string      `json:"meta"`
// 	Level       string      `json:"level"`
// 	Label       string      `json:"label"`
// }

func (l Log) Info(context, message, scope, meta string) {
	if l.logLevel <= 1 {
		_, file, line, _ := runtime.Caller(1)
		msg := fmt.Sprintf("[INFO] Service: %s - Context: %s - Message: %s - Scope: %s - Meta: %s - At: %s:%d", "MESSAGING_SERVICE", context, message, scope, meta, file, line)
		println(msg)
		// json.Unmarshal([]byte(`{"service": "MESSAGING_SERVICE"}`), &serviceName)
		// ls, _ := json.Marshal(logstashPayload{
		// 	ServiceName: serviceName,
		// 	Level:       "info",
		// 	Meta:        meta,
		// 	Message:     message,
		// 	Context:     context,
		// 	Label:       l.appName,
		// 	Scope:       scope,
		// })
		// go l.sendLogstash(string(ls))
	}
}

func (l Log) Error(context, message, scope, meta string) {
	if l.logLevel <= 2 {
		_, file, line, _ := runtime.Caller(1)
		_, file2, line2, _ := runtime.Caller(2)
		msg := fmt.Sprintf("[ERROR] Context: %s - Message: %s - Scope: %s - Meta: %s - At Level1: %s:%d - At Level2: %s:%d", context, message, scope, meta, file, line, file2, line2)
		println(msg)

		// json.Unmarshal([]byte(`{"service": "MESSAGING_SERVICE"}`), &serviceName)
		// ls, _ := json.Marshal(logstashPayload{
		// 	ServiceName: serviceName,
		// 	Level:       "error",
		// 	Meta:        meta,
		// 	Message:     message,
		// 	Context:     context,
		// 	Label:       l.appName,
		// 	Scope:       scope,
		// })
		// go l.sendLogstash(string(ls))
	}
}

func (l Log) Slow(context, message, scope, meta string) {
	if l.logLevel <= 1 {
		_, file, line, _ := runtime.Caller(2)
		msg := fmt.Sprintf("[SLOW] Context: %s - Message: %s - Scope: %s - Meta: %s - At: %s:%d", context, message, scope, meta, file, line)
		println(msg)

		// json.Unmarshal([]byte(`{"service": "MESSAGING_SERVICE"}`), &serviceName)
		// ls, _ := json.Marshal(logstashPayload{
		// 	ServiceName: serviceName,
		// 	Level:       "slow",
		// 	Meta:        meta,
		// 	Message:     message,
		// 	Context:     context,
		// 	Label:       l.appName,
		// 	Scope:       scope,
		// })
		// go l.sendLogstash(string(ls))
	}
}

// func (l Log) sendLogstash(msg string) {

// 	if err := logstsh.Writeln(string(msg)); err != nil {
// 		println(err.Error())
// 	}
// }
