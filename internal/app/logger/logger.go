// The package logger implements logging capabilities in the application.
package logger

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const permFile = 0o644

// AppLogger is a structure with reference to zap.Sugared Logger instance.
type AppLogger struct {
	LZ *zap.SugaredLogger
}

// Write implements io.Writer.
func (a *AppLogger) Write(b []byte) (n int, err error) {
	a.LZ.Errorw(string(b))

	return len(b), nil
}

// NewLogger returns a pointer of struct for app logging.
func NewLogger() *AppLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.StacktraceKey = ""
	fileEncoder := zapcore.NewJSONEncoder(config)
	logFile, err := os.OpenFile("log.json",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, permFile)
	// if an error is written to the console
	if err != nil {
		config := zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.StacktraceKey = ""

		logger, err := config.Build()
		if err != nil {
			log.Fatal(err)
		}

		return &AppLogger{
			LZ: logger.Sugar(),
		}
	}

	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)

	logger := zap.New(core, zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel))

	return &AppLogger{
		LZ: logger.Sugar(),
	}
}
