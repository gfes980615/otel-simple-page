package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	OtelHttpEndpoint   string
	OtelGrpcEndpoint   string
	JaegerEndpoint     string
	Service            string
	MaxQueueSize       int
	MaxExportBatchSize int
	BatchTimeout       int // second
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
		OtelHttpEndpoint:   viper.GetString("config.endpoint.otel_http"),
		OtelGrpcEndpoint:   viper.GetString("config.endpoint.otel_grpc"),
		JaegerEndpoint:     viper.GetString("config.endpoint.jaeger"),
		Service:            viper.GetString("config.service"),
		MaxQueueSize:       viper.GetInt("config.exporter.max_queue_size"),
		MaxExportBatchSize: viper.GetInt("config.exporter.max_export_batch_size"),
		BatchTimeout:       viper.GetInt("config.exporter.batch_timeout"),
	}
}
