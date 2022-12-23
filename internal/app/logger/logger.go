package logger

import "go.uber.org/zap"

type AppLogger struct {
	logger *zap.SugaredLogger
}

func (a *AppLogger) Write(b []byte) (n int, err error) {
	a.logger.Errorw(string(b))
	return len(b), nil
}

func NewLogger() (*AppLogger, error) {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &AppLogger{
		logger: logger.Sugar(),
	}, nil
}
