package logger

import (
	"os"
	"runtime"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	maximumCallerDepth int = 25
	knownLoggerFrames  int = 4
)

// Logger is a struct.
type Logger struct {
	Log      *LogWrapper
	Callback func()
}

// New is a constructor will initialize Logger.
func New(config zap.Config, options ...Option) *Logger {
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	logger, _ := config.Build()
	loggerWithSugar := &Logger{
		Log: &LogWrapper{
			Log:            logger,
			withFileTrace:  sliceContains([]string{"production", "staging", "development"}, os.Getenv(env)),
			withEventTrace: sliceContains([]string{"production", "staging", "development"}, os.Getenv(env)),
			withDDTrace:    sliceContains([]string{"production", "staging", "development"}, os.Getenv(env)),
		},
	}

	for _, opt := range options {
		opt(loggerWithSugar)
	}

	return loggerWithSugar
}

func getCurrentPosition() (fnc string, file string, line int) {
	fr := getCaller(5)

	return fr.Function, fr.File, fr.Line
}

func getCaller(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// sliceContains is a function uses to check slice of strings is contains a string.
func sliceContains(slice []string, str string) bool {
	for _, a := range slice {
		if a == str {
			return true
		}
	}

	return false
}
