package postgre

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	goGS "github.com/dsbasko/go-gs"

	"github.com/dsbasko/pass-keeper/internal/server/config"
	errWrapper "github.com/dsbasko/pass-keeper/pkg/err-wrapper"
	"github.com/dsbasko/pass-keeper/pkg/logger"
)

type Postgre struct {
	db *pgxpool.Pool
}

type Options struct {
	Logger logger.Logger
	Cfg    *config.Config
	GS     goGS.GracefulShutdowner
}

func New(ctx context.Context, opts Options) (_ *Postgre, err error) {
	defer errWrapper.PtrWithOP(&err, "postgre.New")

	// Валидация аргументов
	switch {
	case ctx == nil:
		return nil, ErrMissingContext
	case opts.Logger == nil:
		return nil, ErrMissingLogger
	case opts.Cfg == nil:
		return nil, ErrMissingCfg
	case opts.GS == nil:
		return nil, ErrMissingGS
	}

	// Конфигурация подключения
	pgConfig, err := pgxpool.ParseConfig(opts.Cfg.Provider.Postgre.DSN)
	if err != nil {
		return nil, errWrapper.WithOP(err, "pgxpool.ParseConfig")
	}
	pgConfig.MaxConns = opts.Cfg.Provider.Postgre.MaxConns

	// Создание подключения
	conn, err := pgxpool.NewWithConfig(ctx, pgConfig)
	if err != nil {
		return nil, errWrapper.WithOP(err, "pgxpool.NewWithConfig")
	}

	// Проверка подключения
	if err = conn.Ping(ctx); err != nil {
		return nil, errWrapper.WithOP(err, "conn.Ping")
	}

	// GS отключение
	opts.GS.Subscribe()
	go gsStop(ctx, opts.GS, conn)

	return &Postgre{
		db: conn,
	}, nil
}

func MustNew(ctx context.Context, opts Options) *Postgre {
	resp, err := New(ctx, opts)
	if err != nil {
		panic(err)
	}
	return resp
}

func gsStop(ctx context.Context, gs goGS.GracefulShutdowner, conn *pgxpool.Pool) {
	defer gs.UnsubscribeFn(func() {
		conn.Close()
	})
	<-ctx.Done()
}
