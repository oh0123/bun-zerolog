# bun-zerolog

A query hook for [uptrace/bun](https://github.com/uptrace/bun) that logs with [rs/zerolog](https://github.com/rs/zerolog).

```bash
go get github.com/oh0123/bun-zerolog
```

With the hook added, All errors will be logged at error level. And if you want log everything, please set `zerolog.SetGlobalLevel(zerolog.DebugLevel)`. If `SlowDuration` is defined, only operations that taking longer than the defined duration will be logged.

## Usage

```go
import (
    "os"

    bunzerolog "github.com/oh0123/bun-zerolog"
    "github.com/rs/zerolog"
)

// zerolog.SetGlobalLevel(zerolog.DebugLevel)

db := bun.NewDB()
log := zerolog.New(os.Stdout).With().Logger()
db.AddQueryHook(bunzerolog.NewQueryHook(bunzerolog.QueryHookOptions{
    Logger:         log,
    SlowDuration:   100 * time.Millisecond,
}))
```
