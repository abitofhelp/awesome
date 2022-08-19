package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"time"
)

const (
	kDefaultConnectionTimeout = 1000 * time.Second
	kDefaultGrpcPort          = 20580
	kDefaultHost              = "localhost"
	kDefaultHttpPort          = 18580
)

type AppConfig struct {
	Runtime struct {
		ConnectionTimeOut time.Duration `mapstructure:"CSC_CONNECTION_TIMEOUT"`
		GrpcPort          uint64        `mapstructure:"CSC_GRPC_PORT"`
		Host              string        `mapstructure:"CSC_HOST"`
		HttpPort          uint64        `mapstructure:"CSC_HTTP_PORT"`
	} `mapstructure:"Runtime"`
}

func NewAppConfig(filename string) (*AppConfig, error) {

	if _, err := os.Stat(filename); err == nil {
		// Development mode
		if err := godotenv.Load(filename); err == nil {
			return buildAppConfig()
		} else {
			return nil, fmt.Errorf("\nfailed to load the application configuration file '%s': %w", filename, err)
		}
	} else {
		// Production mode
		err := godotenv.Load()
		if err != nil {
			return nil, err
		}
		return buildAppConfig()
	}
}

func buildAppConfig() (*AppConfig, error) {
	appcfg := &AppConfig{}

	if v, ok := os.LookupEnv("CSC_CONNECTION_TIMEOUT"); ok {
		if val, err := strconv.ParseInt(v, 10, 64); err == nil {
			appcfg.Runtime.ConnectionTimeOut = time.Duration(val) * time.Second
		} else {
			return nil, fmt.Errorf("\nfailed to convert the CSC_GRPC_PORT to a uint value: %w", err)
		}

	} else {
		appcfg.Runtime.ConnectionTimeOut = kDefaultConnectionTimeout
	}

	if v, ok := os.LookupEnv("CSC_GRPC_PORT"); ok {
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			appcfg.Runtime.GrpcPort = val
		} else {
			return nil, fmt.Errorf("\nfailed to convert the CSC_GRPC_PORT to a uint value: %w", err)
		}
	} else {
		appcfg.Runtime.GrpcPort = kDefaultGrpcPort
	}

	if v, ok := os.LookupEnv("CSC_HOST"); ok {
		appcfg.Runtime.Host = v
	} else {
		appcfg.Runtime.Host = kDefaultHost
	}

	if v, ok := os.LookupEnv("CSC_HTTP_PORT"); ok {
		if val, err := strconv.ParseUint(v, 10, 64); err == nil {
			appcfg.Runtime.HttpPort = val
		} else {
			return nil, fmt.Errorf("\nfailed to convert the CSC_HTTP_PORT to a uint value: %w", err)
		}
	} else {
		appcfg.Runtime.HttpPort = kDefaultHttpPort
	}

	return appcfg, nil
}
