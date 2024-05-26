package app

import (
	"context"
	"fmt"
	"github.com/dsbasko/pass-keeper/internal/client/endpoint/tui"
	"github.com/dsbasko/pass-keeper/internal/client/provider"
	"os"
)

func tuiInit(ctx context.Context, cmdCh chan string) *tui.TUI {
	tuiApp := tui.Init(ctx, cmdCh, provider.Init(provider.Mock))
	tuiApp.ExitFn = func() { cmdCh <- tui.Exit }
	tuiApp.LogoutFn = func() { cmdCh <- tui.Logout }
	go func() {
		defer tuiApp.App.EnableMouse(false)
		err := tuiApp.App.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("exit tui app")
		os.Exit(0)
	}()
	return tuiApp
}

func Run(ctx context.Context) (err error) {
	cmdCh := make(chan string)
	tuiApp := tuiInit(ctx, cmdCh)
	defer tuiApp.App.Stop()
	cmdCh <- tui.Vaults //точка входа в TUI, в реале будет зависеть от состояния аутентификации, либо tui.StartMenu, либо tui.EnterPIN
	tuiApp.App.Sync()
	<-ctx.Done()
	return nil
}
