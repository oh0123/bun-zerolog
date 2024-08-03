package bunzerolog_test

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	bunzerolog "github.com/oh0123/bun-zerolog"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
)

func TestQueryHookError(t *testing.T) {
	out := &bytes.Buffer{}
	log := zerolog.New(out).With().Logger()

	h := bunzerolog.NewQueryHook(bunzerolog.QueryHookOptions{
		Logger: log,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM example WHERE id = ?",
		QueryArgs: []interface{}{1},
		Err:       errors.New("database error"),
	}

	h.AfterQuery(context.Background(), event)

	assert.Equal(t, out.String(), `{"level":"error","query":"SELECT * FROM example WHERE id = ?","args":[1],"elapse":0,"error":"database error"}`+"\n")

}

func TestQueryHookDebug(t *testing.T) {
	out := &bytes.Buffer{}

	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	log := zerolog.New(out).With().Logger()

	h := bunzerolog.NewQueryHook(bunzerolog.QueryHookOptions{
		SlowDuration: 100 * time.Millisecond,
		Logger:       log,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM example WHERE id = ?",
		QueryArgs: []interface{}{1},
	}

	h.AfterQuery(context.Background(), event)

	assert.Equal(t, out.String(), `{"level":"debug","query":"SELECT * FROM example WHERE id = ?","args":[1],"elapse":0}`+"\n")

}

func TestQueryHookInfo(t *testing.T) {
	out := &bytes.Buffer{}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	log := zerolog.New(out).With().Logger()

	h := bunzerolog.NewQueryHook(bunzerolog.QueryHookOptions{
		SlowDuration: 100 * time.Millisecond,
		Logger:       log,
	})

	event := &bun.QueryEvent{
		StartTime: time.Now(),
		Query:     "SELECT * FROM example WHERE id = ?",
		QueryArgs: []interface{}{1},
	}

	h.AfterQuery(context.Background(), event)

	assert.Equal(t, out.String(), ``)

}

func TestQueryHookSlow(t *testing.T) {
	out := &bytes.Buffer{}
	log := zerolog.New(out).With().Logger()

	h := bunzerolog.NewQueryHook(bunzerolog.QueryHookOptions{
		SlowDuration: 100 * time.Millisecond,
		Logger:       log,
	})
	event := &bun.QueryEvent{
		StartTime: time.Now().Add(-300 * time.Millisecond),
		Query:     "SELECT * FROM example WHERE id = ?",
		QueryArgs: []interface{}{1},
	}
	h.AfterQuery(context.Background(), event)

	assert.Equal(t, out.String(), `{"level":"warn","query":"SELECT * FROM example WHERE id = ?","args":[1],"elapse":300,"time":"2024-08-03T11:27:11+08:00","message":"slow_query"}`+"\n")
}
