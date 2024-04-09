package app

import (
	"context"

	tui "github.com/dsbasko/pass-keeper/internal/client/endpoint/lazy-tui"
)

func Run(ctx context.Context) (err error) {
	return tui.Run(ctx)
}

func MustRun(ctx context.Context) {
	if err := Run(ctx); err != nil {
		panic(err)
	}
}
