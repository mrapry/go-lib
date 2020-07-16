package logger

import (
	"github.com/brianvoe/gofakeit/v5"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLog(t *testing.T) {
	testCase := map[string]struct {
		level zapcore.Level
	}{
		"Test #1 logger zap log debug": {
			level: zapcore.DebugLevel,
		},
		"Test #2 logger zap log info": {
			level: zapcore.InfoLevel,
		},
		"Test #3 logger zap log warn": {
			level: zapcore.WarnLevel,
		},
		"Test #4 logger zap log error": {
			level: zapcore.ErrorLevel,
		},
		"Test #5 logger zap log panic": {
			level: zapcore.PanicLevel,
		},
	}

	for name, test := range testCase {
		t.Run(name, func(t *testing.T) {
			defer func() {
				recover()
			}()

			Log(test.level, gofakeit.Word(), gofakeit.Word(), gofakeit.Word())
		})
	}
}

func TestLogE(t *testing.T) {
	testName := "Test logger log e"

	t.Run(testName, func(t *testing.T) {
		LogE(gofakeit.Word())
	})
}

func TestLogI(t *testing.T) {
	testName := "Test logger log i"

	t.Run(testName, func(t *testing.T) {
		LogI(gofakeit.Word())
	})
}

func TestLogEf(t *testing.T) {
	testName := "Test logger log ef"

	t.Run(testName, func(t *testing.T) {
		LogEf(gofakeit.Word())
	})
}

func TestLogIf(t *testing.T) {
	testName := "Test logger log if"

	t.Run(testName, func(t *testing.T) {
		LogIf(gofakeit.Word())
	})
}
