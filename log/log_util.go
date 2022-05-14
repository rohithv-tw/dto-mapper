package log

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"time"
)

type LogUtil interface {
	Info(message string)
	Infof(format string, values ...interface{})
	Error(message string)
	Errorf(format string, values ...interface{})
}

type logUtil struct {
	logger *zerolog.Logger
	fields map[string]interface{}
}

func NewLoggerInContext() context.Context {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	childCtx := logger.WithContext(context.Background())

	return childCtx
}

func GetLogger(ctx context.Context, fields map[string]interface{}) LogUtil {
	return &logUtil{
		logger: log.Ctx(ctx),
		fields: fields,
	}
}

func (util *logUtil) Info(message string) {
	util.logger.Info().Fields(util.fields).Msg(message)
}

func (util *logUtil) Infof(format string, values ...interface{}) {
	util.logger.Info().Fields(util.fields).Msgf(format, values)
}

func (util *logUtil) Error(message string) {
	util.logger.Error().Fields(util.fields).Msg(message)
}

func (util *logUtil) Errorf(format string, values ...interface{}) {
	util.logger.Error().Fields(util.fields).Msgf(format, values)
}
