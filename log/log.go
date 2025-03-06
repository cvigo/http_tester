package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var BaseLogger *zap.Logger
var Logger *zap.SugaredLogger

func LogInit(level string, format string) error {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err // and leave the default level
	}
	cfg := zap.Config{
		Level:             l,
		Development:       false,
		DisableCaller:     true,
		Encoding:          format,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		OutputPaths:       []string{"stdout"},
		DisableStacktrace: true,
	}
	if format == "console" {
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	BaseLogger, err = cfg.Build()
	if err != nil {
		return err
	}

	Logger = BaseLogger.Sugar()
	return nil
}

func IsLevelEnabled(level string) bool {
	l, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return false
	}

	return BaseLogger.Core().Enabled(l.Level())
}

func ReportTime(startTs time.Time, msg string, err error) {
	timeNow := time.Now()
	if err != nil {
		Logger.Errorw(err.Error(),
			"start", startTs,
			"finish", timeNow,
			"elapsed(ms)", float64(timeNow.Sub(startTs).Nanoseconds())/1e6,
			"message", err)
	} else {
		Logger.Infow(msg,
			"start", startTs,
			"finish", timeNow,
			"elapsed(ms)", float64(timeNow.Sub(startTs).Nanoseconds())/1e6,
			"message", msg)
	}
}
