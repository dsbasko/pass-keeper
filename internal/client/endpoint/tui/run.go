// Package tui отвечает за отрисовку экранов. Данные и методы получает снаружи.
package tui

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dsbasko/pass-keeper/internal/client/models"
	"github.com/google/uuid"
	"github.com/rivo/tview"
	"os"
)

const (
	minGridWidth  = 100
	minGridHeight = 0
	fieldHeight   = 100
	fieldWidth    = 0
	fieldWidth10  = 10
	fieldWidth20  = 20
	fieldWidth36  = 36
)

type TUI struct {
	email     string
	password  string
	secretKey string
	State     string
	App       *tview.Application
	// fncs      map[string]func() any
	cmdCh     chan string
	ctx       context.Context
	pinRepeat string
	pin       string
	ExitFn    func()
	LogoutFn  func()
	VaultsFn  func() ([]models.VaultForView, error)
	SecretsFn func(string) ([]models.SecretForView, error)
	provider  models.Providerer
}

var t = &TUI{}

func (t *TUI) YesNow(label string, yes, no func()) {
	form := tview.NewForm().
		AddButton("Yes", yes).
		AddButton("No", no)

	form.SetBorder(true).
		SetTitle(fmt.Sprintf("Are you sure you want %s?", label)).
		SetTitleAlign(tview.AlignCenter)
	t.App.QueueUpdateDraw(func() {
		t.App.SetRoot(form, true)
	})
}

func (t *TUI) StartMenuScreen(register, login, quit func()) {
	form := tview.NewForm().
		AddButton("Register", register).
		AddButton("Login", login).
		AddButton("Quit", quit)

	form.SetBorder(true).
		SetTitle("Register or login, please").
		SetTitleAlign(tview.AlignCenter)
	t.App.SetRoot(form, true)
}

func (t *TUI) LoginScreen(login, quit func()) {
	form := tview.NewForm().
		AddInputField("email", "", fieldWidth20, nil, func(text string) {
			t.email = text
		}).
		AddPasswordField("password", "", fieldWidth10, '*', func(text string) {
			t.password = text
		}).
		AddButton("Login", login).
		AddButton("Quit", quit)

	form.SetBorder(true).SetTitle("Login, please").SetTitleAlign(tview.AlignLeft)
	t.App.SetRoot(form, true)
}

func (t *TUI) RegisterScreen(register, quit func()) {
	form := tview.NewForm().
		AddInputField("email", "", fieldWidth20, nil, func(text string) {
			t.email = text
		}).
		AddPasswordField("password", "", fieldWidth20, '*', func(text string) {
			t.password = text
		}).
		AddInputField("secretKey", uuid.New().String(), fieldWidth36, nil, func(text string) {
			t.secretKey = text
		}).
		AddButton("Register", register).
		AddButton("Quit", quit)

	form.SetBorder(true).SetTitle("Register, please").SetTitleAlign(tview.AlignLeft)
	t.App.SetRoot(form, true)
}

func (t *TUI) CreatePINScreen(regPin, quit func()) {
	form := tview.NewForm().
		AddPasswordField("PIN (10 symbols)", "", fieldWidth10, '*', func(text string) {
			t.pin = text
		}).
		AddPasswordField("repeat PIN", "", fieldWidth10, '*', func(text string) {
			t.pinRepeat = text
		}).
		AddButton("Set PIN", regPin).
		AddButton("Quit", quit)

	form.SetBorder(true).SetTitle("Register, please").SetTitleAlign(tview.AlignLeft)
	t.App.SetRoot(form, true)
}

func (t *TUI) EnterPINScreen(enterPIN, quit func()) {
	form := tview.NewForm().
		AddPasswordField("PIN (10 symbols)", "", fieldWidth10, '*', func(text string) {
			t.pin = text
		}).
		AddButton("enter PIN", enterPIN).
		AddButton("Quit", quit)

	form.SetBorder(true).SetTitle("Register, please").SetTitleAlign(tview.AlignLeft)
	t.App.SetRoot(form, true)
}

func (t *TUI) MainMenuScreen(vaults, logout, quit func()) {
	form := tview.NewForm().
		AddButton("Vaults", vaults).
		AddButton("Logout", logout).
		AddButton("Exit", quit)

	form.SetBorder(true).
		SetTitle("Register or login, please").
		SetTitleAlign(tview.AlignCenter)
	t.App.SetRoot(form, true)
}

func (t *TUI) VaultsScreen() {
	const rows = 5
	grid := tview.NewGrid().
		SetRows(rows, 0).
		SetColumns(-2, -2, -1).
		SetBorders(true)

	secretScreen := tview.NewForm()
	secretScreen.SetTitle("Secret")

	secretsTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(false, false)
	secretsTable.SetTitle("Secrets")

	vaultsTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(false, false)
	vaultsTable.SetTitle("Vaults")
	vaultsTable.
		SetCellSimple(0, 0, "Name").
		SetCellSimple(0, 1, "Data")
	vaults, err := t.provider.GetVaults()
	if err != nil {
		//TODO модальное окно с ошибками
	}
	for k, v := range vaults {
		cellName := tview.NewTableCell(v.Name).
			SetSelectable(true).
			SetClickedFunc(func() bool {
				secretsTable.Clear()
				secretsTable.
					SetCellSimple(0, 0, "Name").
					SetCellSimple(0, 1, "Type").
					SetCellSimple(0, 2, "Data")
				secrets, err := t.provider.GetSecrets(v.ID)
				if err != nil {
					//TODO модальное окно с ошибками
				}
				for k, secret := range secrets {
					cellName := tview.NewTableCell(secret.Name).
						SetClickedFunc(func() bool {
							dataSecret, _ := json.MarshalIndent(secret, "", "  ")

							secretScreen.Clear(true)
							secretScreen.
								AddButton("view", func() {}).
								AddButton("edit", func() {}).
								AddTextView("Secret", string(dataSecret), fieldWidth, fieldHeight, true, true)

							grid.
								AddItem(secretScreen, 1, 2, 1, 1, minGridHeight, minGridWidth, true)

							return true
						})
					cellType := tview.NewTableCell(secret.Type)
					cellData := tview.NewTableCell(fmt.Sprintf("[%s][%s]\n[%s]",
						secret.VaultID,
						v.CreateAt.Format("2006-01-02 15:04:05"),
						v.Comment,
					)).SetExpansion(1)
					secretsTable.SetCell(k+1, 0, cellName)
					secretsTable.SetCell(k+1, 1, cellType)
					secretsTable.SetCell(k+1, 2, cellData)
				}

				grid.AddItem(vaultsTable, 1, 0, 1, 1, minGridHeight, minGridWidth, true)
				grid.AddItem(secretsTable, 1, 1, 1, 1, minGridHeight, minGridWidth, false)
				t.App.Sync()
				return true
			})
		cellComment := tview.NewTableCell(fmt.Sprintf("[%s]\n[%s]",
			v.CreateAt.Format("2006-01-02 15:04:05"),
			v.Comment,
		)).SetExpansion(1)

		vaultsTable.
			SetCell(k+1, 0, cellName).
			SetCell(k+1, 1, cellComment)
	}

	menu := tview.NewForm().
		AddButton("Logout", t.LogoutFn).
		AddButton("Exit", t.ExitFn)

	menu.SetBorder(true).
		SetTitle("Menu").
		SetTitleAlign(tview.AlignCenter)

	grid.
		AddItem(vaultsTable, 1, 0, 1, 1, minGridHeight, minGridWidth, true).
		AddItem(secretsTable, 1, 1, 1, 1, minGridHeight, minGridWidth, false).
		AddItem(secretScreen, 1, 2, 1, 1, minGridHeight, minGridWidth, false).
		AddItem(menu, 0, 0, 1, 3, minGridHeight, 0, false)

	t.App.SetRoot(grid, true)
}
func (t *TUI) VaultEditScreen()  {}
func (t *TUI) SecretViewScreen() {}
func (t *TUI) SecretEditScreen() {}

func (t *TUI) mainLoop() {
	defer t.App.EnableMouse(false)
	for {
		select {
		case <-t.ctx.Done():
			t.App.Stop()
			t = nil
			return
		case cmd := <-t.cmdCh:
			switch cmd {
			// Я думал над тем, чтобы сделать мапу,
			// но это была бы мапа map[string]interface{}
			// потому что функции с разными сигнатурами,
			// не получится их хранить и не кастовать
			// так что смысла в этом нет(ИМХО)
			case Exit:
				t.YesNow(Exit, func() {
					t.App.EnableMouse(false)
					os.Exit(0)
				},
					func() { t.cmdCh <- MainMenu })
			case Logout:
				t.YesNow(Logout, func() { os.Exit(0) },
					func() { t.cmdCh <- MainMenu })
			case StartMenu:
				t.StartMenuScreen(func() { t.cmdCh <- RegisterMenu },
					func() { t.cmdCh <- LoginMenu }, exitApp(t.cmdCh))
			case LoginMenu:
				t.LoginScreen(func() {
					t.cmdCh <- EnterPIN
				}, exitApp(t.cmdCh))
			case RegisterMenu:
				t.RegisterScreen(func() {
					t.cmdCh <- CreatePIN
				}, exitApp(t.cmdCh))
			case CreatePIN:
				t.CreatePINScreen(func() {
					t.cmdCh <- LoginMenu
				}, exitApp(t.cmdCh))
			case EnterPIN:
				t.EnterPINScreen(func() {
					t.cmdCh <- Vaults
				}, exitApp(t.cmdCh))
			case Vaults:
				t.VaultsScreen()
			case VaultEdit:
				t.VaultEditScreen()
			case SecretView:
				t.SecretViewScreen()
			case SecretEdit:
				t.SecretEditScreen()
			}
		}
	}
}

func exitApp(cmdCh chan string) func() {
	return func() {
		cmdCh <- Exit
	}
}

func Init(ctx context.Context, ch chan string, provider models.Providerer) *TUI {
	t.App = tview.NewApplication()
	t.App.EnableMouse(true)
	t.cmdCh = ch
	t.ctx = ctx
	t.provider = provider

	go t.mainLoop()

	return t
}
