package config

import (
	"github.com/SiriusServiceDesk/notification-service/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"os"
)

type GrpcServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"3000"`
}

type HttpServer struct {
	Host string `yaml:"host" env-default:"0.0.0.0"`
	Port string `yaml:"port" env-default:"8000"`
}

type Db struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"postgres"`
	Name     string `yaml:"db_name" env-default:"postgres"`
}

type RabbitMq struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     string `yaml:"port" env-default:"6379"`
	User     string `yaml:"user" env-default:"guest"`
	Password string `yaml:"password" env-default:"guest"`
}

type Email struct {
	SmtpHost     string `yaml:"smtp_host"`
	SmtpPort     int    `yaml:"smtp_port"`
	SmtpUser     string `yaml:"smtp_user"`
	SmtpPassword string `yaml:"smtp_password"`
}

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	HttpServer HttpServer `yaml:"http_server"`
	GrpcServer GrpcServer `yaml:"grpc_server"`
	Db         Db         `yaml:"db"`
	RabbitMQ   RabbitMq   `yaml:"rabbit"`
	Email      Email      `yaml:"email"`
}

func GetConfig() *Config {

	path := os.Getenv("CONFIG_PATH")
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file does not exist " + path)
	}

	var config Config

	if err := cleanenv.ReadConfig(path, &config); err != nil {
		logger.Fatal("cannot read config file", zap.String("path", path), zap.Error(err))
	}

	return &config
}
