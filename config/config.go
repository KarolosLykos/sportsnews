package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Dev      bool `envconfig:"DEV" default:"true"`
	HTTP     HTTP
	Logger   Logger
	MongoDB  MongoConfig
	Redis    RedisConfig
	Consumer ConsumerConfig
}

type HTTP struct {
	Port string `envconfig:"HTTP_PORT" default:":8081"`
}

type Logger struct {
	LogLevel string `envconfig:"LOGGER_LEVEL" default:"debug"`
	Format   string `envconfig:"LOGGER_FORMAT" default:"json"`
}

type MongoConfig struct {
	Host     string `envconfig:"MONGO_HOST" default:"localhost"`
	Port     string `envconfig:"MONGO_PORT" default:"27017"`
	Username string `envconfig:"MONGO_Username" default:"admin"`
	Password string `envconfig:"MONGO_Password" default:"secret"`
}

type RedisConfig struct {
	Host       string        `envconfig:"REDIS_HOST" default:"localhost"`
	Port       string        `envconfig:"REDIS_PORT" default:"6379"`
	Expiration time.Duration `envconfig:"REDIS_EXPIRATION" default:"3600s"`
	KeyPrefix  string        `envconfig:"REDIS_KEY_PREFIX" default:"articles"`
}

type ConsumerConfig struct {
	HullConsumer HullConsumer
}

type HullConsumer struct {
	Frequency time.Duration `envconfig:"HULL_CONSUMER_FREQUENCY" default:"30m"`
	SingleURL string        `envconfig:"HULL_CONSUMER_SINGLE_URL" default:"https://www.wearehullcity.co.uk/api/incrowd/getnewsarticleinformation"`
	ListURL   string        `envconfig:"HULL_CONSUMER_LIST_URL" default:"https://www.wearehullcity.co.uk/api/incrowd/getnewlistinformation"`
	Count     int           `envconfig:"HULL_CONSUMER_COUNT" default:"50"`
}

func Parse() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("could not read env file: %v", err)
	}

	return cfg, nil
}
