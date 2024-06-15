package logger

import (
	"context"
	"github.com/jackc/pgx/v5/tracelog"
	"log"
)

type Logger struct {
}

func (l Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]any) {
	log.Println(msg, data)
}
