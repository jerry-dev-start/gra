package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	encoderCfg.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	}

	core := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level)
	Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()

	// 接管 Gin 内部 debug 日志输出
	gin.DefaultWriter = newZapWriter(Log, zapcore.InfoLevel)
	gin.DefaultErrorWriter = newZapWriter(Log, zapcore.ErrorLevel)
}

// GinLogger 替代 gin 默认日志中间件，统一走 zap
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if query != "" {
			path = path + "?" + query
		}

		Log.Infof("%d | %13v | %15s | %-7s %s",
			status, latency, clientIP, method, path,
		)

		// 记录错误
		if len(c.Errors) > 0 {
			Log.Errorf("gin errors: %s", c.Errors.ByType(gin.ErrorTypePrivate).String())
		}
	}
}

// GinRecovery 替代 gin 默认 recovery 中间件，panic 走 zap
func GinRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		Log.Errorf("panic recovered: %v | %s %s", recovered, c.Request.Method, c.Request.URL.Path)
		c.AbortWithStatus(500)
	})
}

// zapWriter 将 io.Writer 接口桥接到 zap，用于接管 gin.DefaultWriter
type zapWriter struct {
	logger *zap.SugaredLogger
	level  zapcore.Level
}

func newZapWriter(logger *zap.SugaredLogger, level zapcore.Level) io.Writer {
	return &zapWriter{logger: logger, level: level}
}

func (w *zapWriter) Write(p []byte) (n int, err error) {
	msg := strings.TrimRight(string(p), "\n")
	switch w.level {
	case zapcore.ErrorLevel:
		w.logger.Error(msg)
	default:
		w.logger.Info(msg)
	}
	return len(p), nil
}
