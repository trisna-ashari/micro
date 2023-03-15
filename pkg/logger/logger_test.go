package logger_test

import (
	"errors"
	"go.uber.org/zap/zapcore"
	"micro/pkg/logger"
	"runtime"
	"testing"

	"go.uber.org/zap"
)

func TestLoggerNew(t *testing.T) {
	logStd := logger.New(zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "file",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.RFC3339TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	})

	runtime.GOMAXPROCS(runtime.NumCPU())

	ch := make(chan string)

	const numReqs = 100
	routine := func() {
		for i := 0; i < numReqs; i++ {
			ch <- logStd.Log.
				WithEvent(&logger.LogEvent{
					Name:   "TestLoggerNew",
					Type:   "unit_test",
					Method: "TestLoggerNew",
					Path:   "/",
				}).
				Test("successfully run logger test, err: %v", errors.New("no error"))
		}
	}

	const numRoutines = 10
	for i := 0; i < numRoutines; i++ {
		go func() {
			routine()
		}()
	}

	total := uint32(numReqs * numRoutines)
	for i := uint32(0); i < total; i++ {
		_ = <-ch
	}
}
