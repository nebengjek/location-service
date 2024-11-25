// bin/app/main_test.go
package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock dependencies
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) GetClient() *MockRedisClient {
	return m
}

type MockKafkaProducer struct {
	mock.Mock
}

func (m *MockKafkaProducer) NewProducer(config map[string]interface{}, logger interface{}) (*MockKafkaProducer, error) {
	return m, nil
}

type MockMongoDBLogger struct {
	mock.Mock
}

func (m *MockMongoDBLogger) NewMongoDBLogger(conn, dbName string, logger interface{}) *MockMongoDBLogger {
	return m
}

func TestMainFunction(t *testing.T) {

	// Initialize Echo server
	e := echo.New()
	configureHttp(e)

	// Create a test server
	server := httptest.NewServer(e)
	defer server.Close()

	// Test health-check endpoint
	resp, err := http.Get(server.URL + "/v1/health-check")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test server shutdown
	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		server.Close()
		<-ctx.Done()
		close(done)
	}()

	quit <- os.Interrupt
	<-done
}

func configureHttp(e *echo.Echo) {
	e.GET("/v1/health-check", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "This service is running properly"})
	})
}
