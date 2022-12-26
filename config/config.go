package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OtelHttpEndpoint string
	OtelGrpcEndpoint string
	JaegerEndpoint   string
	Service          string
}

var AppConfig *Config

func init() {
	viper.SetConfigName("app.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	fmt.Println(viper.GetString("config.endpoint.otel_http"))
	AppConfig = &Config{
		OtelHttpEndpoint: viper.GetString("config.endpoint.otel_http"),
		OtelGrpcEndpoint: viper.GetString("config.endpoint.otel_grpc"),
		JaegerEndpoint:   viper.GetString("config.endpoint.jaeger"),
		Service:          viper.GetString("config.service"),
	}
}
