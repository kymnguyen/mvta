package config

import "time"

type Config struct {
	AppEnv string
	HTTP   HTTPConfig
	Mongo  MongoConfig
}

type HTTPConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type MongoConfig struct {
	URI      string
	Database string
}
