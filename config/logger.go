package config

import (
	"context"
	"os"

	"go.uber.org/zap/zapcore"

	zapLog "go.uber.org/zap"
)

type correlationIdType int
const (
	requestIdKey correlationIdType = iota
)

type loggerKeyType int

const loggerKey loggerKeyType = iota

var sugarLogger *zapLog.SugaredLogger

func getEncoder() zapcore.Encoder {
	encoderConfig := zapLog.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func InitLogger() *zapLog.SugaredLogger {

	swSugar := zapcore.AddSync(os.Stdout)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, swSugar, zapcore.DebugLevel)

	zapLogger := zapLog.New(core, zapLog.AddCaller())
	zapLog.ReplaceGlobals(zapLogger)

	sugarLogger = zapLogger.Sugar()

	return sugarLogger
}

func NewContext(ctx context.Context, fields interface{}) context.Context {
	return context.WithValue(ctx, loggerKey, WithContext(ctx).With(fields))
}

func WithContext(ctx context.Context) *zapLog.SugaredLogger {
	if ctx == nil {
		return sugarLogger
	}

	if ctxLogger, ok := ctx.Value(loggerKey).(*zapLog.SugaredLogger); ok {
		return ctxLogger
	}

	return sugarLogger
}

// WithRqId returns a context which knows its request ID
func WithRqId(ctx context.Context, rqId string) context.Context {
	return context.WithValue(ctx, requestIdKey, rqId)
}

// Logger returns a zap logger with as much context as possible
func Logger(ctx context.Context) *zapLog.SugaredLogger {
	newLogger := sugarLogger
	if ctx != nil {
		if ctxRqId, ok := ctx.Value(requestIdKey).(string); ok {
			newLogger = newLogger.With(zapLog.String("cID", ctxRqId))
		}
	}
	return newLogger
}
