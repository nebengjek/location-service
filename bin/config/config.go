package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type envConfig struct {
	APMSecretToken    string
	APMUrl            string
	AppEnv            string
	AppName           string
	AppPort           string
	AppVersion        string
	BasicAuthPassword string
	BasicAuthUsername string
	CipherKey         string
	ConfigCors        string
	IvKey             string
	KafkaUrl          string
	LogLevel          string
	LogstashHost      string
	LogstashPort      string
	MinioAccessKey    string
	MinioEndpoint     string
	MinioSecretKey    string
	MinioUseSSL       bool
	MongoMasterDBUrl  string
	MongoSlaveDBUrl   string
	PrivateKey        string
	PublicKey         string
	JwtAudience       string
	JwtIssuer         string
	JwtAlgorithm      string
	ShutdownDelay     int

	RedisHost            string
	RedisPort            string
	RedisPassword        string
	RedisDB              string
	RedisClusterNode     string
	RedisClusterPassword string
	UseRedis             bool
	ElasticHost          string
	ElasticUsername      string
	ElasticPassword      string
	ElasticMaxRetries    int
	GoogleApiKey         string
	SocketUrl            string
}

func (e envConfig) LogstashPortInt() int {
	i, err := strconv.ParseInt(e.LogstashPort, 10, 64)
	if err != nil {
		panic(err)
	}

	return int(i)
}

func (e envConfig) DnsMariaDB() (string, string) {
	var (
		mariaDbHost     = os.Getenv("MYSQL_HOST")
		mariaDbUsername = os.Getenv("MYSQL_USERNAME")
		mariaDbPassword = os.Getenv("MYSQL_PASSWORD")
		mariaDbName     = os.Getenv("MYSQL_DB_NAME")
	)
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", mariaDbUsername, mariaDbPassword, mariaDbHost, mariaDbName), mariaDbName

}

var envCfg envConfig

func init() {
	err := godotenv.Load()

	if err != nil {
		println(err.Error())
	}

	shutdownDelay, _ := strconv.Atoi(os.Getenv("SHUTDOWN_DELAY"))                // default 0
	minioUseSsl, _ := strconv.ParseBool(os.Getenv("MINIO_USE_SSL"))              // default false
	UseRedis, _ := strconv.ParseBool(os.Getenv("REDIS_CONFIG_CLUSTER"))          // default false
	elasticMaxRetries, _ := strconv.Atoi(os.Getenv("ELASTICSEARCH_MAX_RETRIES")) // default false

	envCfg = envConfig{
		APMSecretToken:       os.Getenv("ELASTIC_APM_SECRET_TOKEN"),
		APMUrl:               os.Getenv("ELASTIC_APM_SERVER_URL"),
		AppEnv:               os.Getenv("APP_ENV"),
		AppName:              os.Getenv("APP_NAME"),
		AppPort:              os.Getenv("APP_PORT"),
		AppVersion:           os.Getenv("APP_VERSION"),
		BasicAuthPassword:    os.Getenv("BASIC_AUTH_PASSWORD"),
		BasicAuthUsername:    os.Getenv("BASIC_AUTH_USERNAME"),
		CipherKey:            os.Getenv("AES_KEY"),
		ConfigCors:           os.Getenv("CORS_CONFIG"),
		IvKey:                "",
		KafkaUrl:             os.Getenv("KAFKA_HOST_URL"),
		LogLevel:             os.Getenv("LOG_LEVEL"),
		LogstashHost:         os.Getenv("LOGSTASH_HOST"),
		LogstashPort:         os.Getenv("LOGSTASH_PORT"),
		MinioAccessKey:       os.Getenv("MINIO_ACCESS_KEY"),
		MinioEndpoint:        os.Getenv("MINIO_END_POINT"),
		MinioSecretKey:       os.Getenv("MINIO_SECRET_KEY"),
		MinioUseSSL:          minioUseSsl,
		MongoMasterDBUrl:     os.Getenv("MONGO_MASTER_DATABASE_URL"),
		MongoSlaveDBUrl:      os.Getenv("MONGO_SLAVE_DATABASE_URL"),
		PrivateKey:           os.Getenv("PRIVATE_KEY_PATH"),
		PublicKey:            os.Getenv("PUBLIC_KEY_PATH"),
		JwtAudience:          os.Getenv("JWT_AUDIENCE"),
		JwtIssuer:            os.Getenv("JWT_ISSUER"),
		JwtAlgorithm:         os.Getenv("JWT_SIGNING_ALGORITHM"),
		ShutdownDelay:        shutdownDelay,
		RedisHost:            os.Getenv("REDIS_HOST"),
		RedisPort:            os.Getenv("REDIS_PORT"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
		RedisDB:              os.Getenv("REDIS_DB"),
		RedisClusterNode:     os.Getenv("REDIS_CLUSTER_NODE"),
		RedisClusterPassword: os.Getenv("REDIS_CLUSTER_PASSWORD"),
		UseRedis:             UseRedis,

		ElasticHost:       os.Getenv("ELASTICSEARCH_HOST"),
		ElasticUsername:   os.Getenv("ELASTICSEARCH_USERNAME"),
		ElasticPassword:   os.Getenv("ELASTICSEARCH_PASSWORD"),
		ElasticMaxRetries: elasticMaxRetries,

		GoogleApiKey: os.Getenv("GOOGLE_API_KEY"),
		SocketUrl:    os.Getenv("SOCKET_URL"),
	}
}

func GetConfig() *envConfig {
	return &envCfg
}
