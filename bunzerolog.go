package bunzerolog

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	"github.com/uptrace/bun"
)

const (
	SqlFieldName            = "query"
	SqlArgsFieldName        = "args"
	SqlExecuteTimeFieldName = "elapse"
)

type QueryHook struct {
	bun.QueryHook
	logger       zerolog.Logger
	slowDuration time.Duration
}

type QueryHookOptions struct {
	Logger       zerolog.Logger
	SlowDuration time.Duration
}

func NewQueryHook(options QueryHookOptions) QueryHook {
	return QueryHook{
		logger:       options.Logger,
		slowDuration: options.SlowDuration,
	}
}

func (h QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}

func (h QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	execElapse := time.Since(event.StartTime)

	sqlLogger := h.logger.With().
		Str(SqlFieldName, event.Query).
		Any(SqlArgsFieldName, event.QueryArgs).
		Int64(SqlExecuteTimeFieldName, execElapse.Milliseconds()).
		Logger()

	if event.Err != nil {
		sqlLogger.Error().Err(event.Err).Msg("")
		return
	}

	if execElapse >= h.slowDuration {
		sqlLogger.Warn().Msg("slow_query")
		return
	}

	sqlLogger.Debug().Msg("")
}
