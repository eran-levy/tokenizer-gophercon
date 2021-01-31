package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

var Log *zap.SugaredLogger

type Config struct {
	LogLevel      string
	ApplicationId string
}

func New(c Config) {
	const serviceKey = "service"
	atom := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	atom.SetLevel(zap.InfoLevel)
	logLevel := c.LogLevel
	if logLevel != "" {
		err := atom.UnmarshalText([]byte(strings.ToLower(logLevel)))
		if err != nil {
			logger.Fatal("Invalid Log level")
		}
	}

	Log = logger.Sugar().With(serviceKey, c.ApplicationId)
}

func Close() {
	_ = Log.Sync()
}
