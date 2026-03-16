package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gra/pkg/config"
)

var Log *zap.SugaredLogger

func Init(cfg *config.LogConfig) {
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.Level))

	var encoder zapcore.Encoder
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
}
