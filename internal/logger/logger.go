package logger

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(level zerolog.Level) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	zerolog.TimeFieldFormat = time.RFC3339Nano + "+09:00"
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	eventsf(ctx, log.Debug(), format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	eventsf(ctx, log.Info(), format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	eventsf(ctx, log.Warn(), format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	eventsf(ctx, log.Error(), format, args...)
}

func Eventsf(ctx context.Context, events *zerolog.Event, format string, args ...interface{}) {
	eventsf(ctx, events, format, args...)
}

func eventsf(ctx context.Context, events *zerolog.Event, format string, args ...interface{}) {
	if len(args) > 0 {
		events.Msgf(format, args...)
	} else {
		events.Msgf(format)
	}
}
