package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/dsbasko/pass-keeper/internal/client/endpoint/tui"
)

func secretsExample(id string) []tui.SecretForView {
	fmt.Println(id)
	TestList := []tui.SecretForView{{
		ID:             "1",
		VaultID:        id,
		Name:           "xxxxxx",
		Secret:         []byte("1243345"),
		Comment:        "some comment",
		UnPackedSecret: "543143221",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "FILE",
	}, {
		ID:             "2",
		VaultID:        id,
		Name:           "zzzzzzz",
		Secret:         []byte("12315tr4"),
		Comment:        "some comment",
		UnPackedSecret: "12312312",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "STRING",
	}, {
		ID:             "3",
		VaultID:        id,
		Name:           "wwwwwww",
		Secret:         []byte("f423g4342"),
		Comment:        "some comment",
		UnPackedSecret: "d1q2wt2t2",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "KEYPASS",
	}, {
		ID:             "4",
		VaultID:        id,
		Name:           "rrrrr",
		Secret:         []byte("f23g43gf2"),
		Comment:        "some comment",
		UnPackedSecret: "d1d12d1221",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "STRING",
	}, {
		ID:             "5",
		VaultID:        id,
		Name:           "qqqqqqqq",
		Secret:         []byte("e12cfw"),
		Comment:        "some comment",
		UnPackedSecret: "d12e12121t34",
		CreateAt:       time.Now(),
		UpdateAt:       time.Now(),
		Type:           "CARD",
	}}
	return TestList
}

// todo удалить
func vaultsExample() []tui.VaultForView {
	var vaults []tui.VaultForView
	for i := 0; i < 10; i++ {
		vaults = append(vaults, tui.VaultForView{
			ID:       fmt.Sprintf("%d", i),
			Name:     fmt.Sprintf("someVault%d", i),
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
			Comment:  "i'm comment, mlem",
		})
	}
	return vaults
}

func exit(cmdCh chan string) func() {
	return func() {
		cmdCh <- tui.Exit
	}
}
func tuiInit(ctx context.Context, cmdCh chan string) *tui.OurTui {
	tuiApp := tui.New(ctx, cmdCh)
	tuiApp.ExitFn = func() { cmdCh <- tui.Exit }
	tuiApp.LogoutFn = func() { cmdCh <- tui.Logout }
	tuiApp.VaultsFn = vaultsExample
	tuiApp.SecretsFn = secretsExample
	go func() {
		defer tuiApp.App.EnableMouse(false)
		err := tuiApp.App.Run()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("exit tui app")
	}()
	return tuiApp
}

func Run(ctx context.Context) (err error) {
	cmdCh := make(chan string)
	tuiApp := tuiInit(ctx, cmdCh)

	go func(commandCh chan string) {
		for {
			select {
			case <-ctx.Done():
				tuiApp.App.Stop()
				tuiApp = nil
				return
			case cmd := <-commandCh:
				switch cmd {
				// Я думал над тем, чтобы сделать мапу,
				// но это была бы мапа map[string]interface{}
				// потому что функции с разными сигнатурами,
				// не получится их хранить и не кастовать
				// так что смысла в этом нет(ИМХО)
				case tui.Exit:
					tuiApp.YesNow(tui.Exit, func() { os.Exit(0) },
						func() { cmdCh <- tui.MainMenu })
				case tui.Logout:
					tuiApp.YesNow(tui.Logout, func() { os.Exit(0) },
						func() { cmdCh <- tui.MainMenu })
				case tui.StartMenu:
					tuiApp.StartMenuScreen(func() { cmdCh <- tui.RegisterMenu },
						func() { cmdCh <- tui.LoginMenu }, exit(cmdCh))
				case tui.LoginMenu:
					tuiApp.LoginScreen(func() {
						cmdCh <- tui.EnterPIN
					}, exit(cmdCh))
				case tui.RegisterMenu:
					tuiApp.RegisterScreen(func() {
						cmdCh <- tui.CreatePIN
					}, exit(cmdCh))
				case tui.CreatePIN:
					tuiApp.CreatePINScreen(func() {
						cmdCh <- tui.LoginMenu
					}, exit(cmdCh))
				case tui.EnterPIN:
					tuiApp.EnterPINScreen(func() {
						cmdCh <- tui.Vaults
					}, exit(cmdCh))
				case tui.Vaults:
					tuiApp.VaultsScreen()
				case tui.VaultEdit:
					tuiApp.VaultEditScreen()
				case tui.SecretView:
					tuiApp.SecretViewScreen()
				case tui.SecretEdit:
					tuiApp.SecretEditScreen()
				}
			}
		}
	}(cmdCh)

	cmdCh <- tui.Vaults
	<-ctx.Done()
	return nil
}
