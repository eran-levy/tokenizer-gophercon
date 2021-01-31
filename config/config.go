package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"strings"
	"time"
)

type ServiceParams struct {
	LogLevel         string
	AppId            string
	DebugModeEnabled bool
}
type RESTApiAdapterParams struct {
	HttpAddress          string
	TerminationTimeout   time.Duration
	ReadRequestTimeout   time.Duration
	WriteResponseTimeout time.Duration
}
type TelemetryParams struct {
	TracingAgentEndpoint string
}
type CacheParams struct {
	CacheSize      int
	CacheAddress   string
	ReadTimeout    time.Duration
	ExpirationTime time.Duration
}
type EnvConfiguration struct {
	Service        ServiceParams
	RESTApiAdapter RESTApiAdapterParams
	Telemetry      TelemetryParams
	Cache          CacheParams
}

// Load default configs. Using spf13/viper here, you can use any library that fit your needs.
// Make sure all defaults are available and override only the configuration that varies between environments.
// Other libraties that can help you and use tags:  kelseyhightower/envconfig, ardanlabs/conf - for example:
// var cfg struct {
//	 LogLevel string `envconfig:"LOG_LEVEL" default:"DEBUG"`
// }
// err := envconfig.Process("", &cfg)
func LoadConfig() (EnvConfiguration, error) {
	const (
		defaultValuesPath = "helm"
		valuesFileName    = "values"
		fileType          = "yaml"
	)
	env := viper.New()
	env.AddConfigPath(defaultValuesPath)
	env.SetConfigName(valuesFileName)
	env.SetConfigType(fileType)
	env.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	env.AutomaticEnv()
	err := env.ReadInConfig() // Find and read the config file
	if err != nil {           // Handle errors reading the config file
		return EnvConfiguration{}, errors.Wrap(err, "could not read env config values")
	}
	return EnvConfiguration{Service: ServiceParams{
		LogLevel:         env.GetString("service.LOG_LEVEL"),
		AppId:            env.GetString("service.APPLICATION_ID"),
		DebugModeEnabled: env.GetBool("service.DEBUG_MODE_ENABLED")},
		RESTApiAdapter: RESTApiAdapterParams{
			HttpAddress:          env.GetString("service.REST_API_HTTP_ADDRESS"),
			TerminationTimeout:   env.GetDuration("service.REST_API_HTTP_TERMINATION_TIMEOUT"),
			ReadRequestTimeout:   env.GetDuration("service.REST_API_HTTP_READ_REQUEST_TIMEOUT"),
			WriteResponseTimeout: env.GetDuration("service.REST_API_HTTP_WRITE_RESPONSE_TIMEOUT"),
		},
		Telemetry: TelemetryParams{TracingAgentEndpoint: env.GetString("service.TELEMETRY_JAEGER_ENDPOINT")},
		Cache: CacheParams{CacheSize: env.GetInt("service.CACHE_SIZE_INT"),
			CacheAddress:   env.GetString("service.CACHE_ADDRESS"),
			ReadTimeout:    env.GetDuration("service.CACHE_READ_TIMEOUT"),
			ExpirationTime: env.GetDuration("service.CACHE_KEYS_EXPIRY_TTL")}}, nil
}
