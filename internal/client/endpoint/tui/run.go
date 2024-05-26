// Package tui отвечает за отрисовку экранов. Данные и методы получает снаружи.
package tui

import (
	"context"
	"fmt"
	"github.com/dsbasko/pass-keeper/internal/client/models"
	"github.com/gdamore/tcell/v2"
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

func newPrimitive(text string) tview.Primitive {
	return tview.NewTextView().
		SetTextAlign(tview.AlignCenter).
		SetText(text)
}

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
	vaultsTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(true, false)
	vaultsTable.SetTitle("Vaults")
	vaultsTable.SetSelectedFunc(func(row, columnt int) {
		cell := vaultsTable.GetCell(row, 0)
		t.SecretsScreen(cell.Text)
	})
	vaultsTable.
		SetCell(0, 0, tview.NewTableCell("ID").SetSelectable(false)).
		SetCell(0, 1, tview.NewTableCell("Name").SetSelectable(false)).
		SetCell(0, 2, tview.NewTableCell("Data").SetSelectable(false))

	vaults, err := t.provider.GetVaults()
	if err != nil {
		//TODO модальное окно с ошибками
	}
	for k, v := range vaults {
		cellID := tview.NewTableCell(v.ID).
			SetClickedFunc(func() bool {
				t.SecretsScreen(v.ID)
				return true
			})
		cellName := tview.NewTableCell(v.Name).
			SetClickedFunc(func() bool {
				t.SecretsScreen(v.ID)
				return true
			})
		cellComment := tview.NewTableCell(fmt.Sprintf("[%s]\n[%s]",
			v.CreateAt.Format("2006-01-02 15:04:05"),
			v.Comment,
		)).SetExpansion(1)

		vaultsTable.
			SetCell(k+1, 0, cellID).
			SetCell(k+1, 1, cellName).
			SetCell(k+1, 2, cellComment)
	}

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(newPrimitive("Vaults"), 0, 10, false)
	flex.AddItem(vaultsTable, 0, 100, true)
	flex.AddItem(newPrimitive("Esc for return, 'e' for edit, 'enter' for view secrets"), 0, 10, false)
	t.App.SetRoot(flex, true)
}

func (t *TUI) SecretsScreen(vaultID string) {
	secretsTable := tview.NewTable().
		SetBorders(true).
		SetFixed(1, 1).
		SetSelectable(true, false)
	secretsTable.SetSelectedFunc(func(row, columnt int) {
		cell := secretsTable.GetCell(row, 0)
		t.SecretViewScreen(cell.Text)
	})
	secretsTable.SetTitle("Secrets")
	secretsTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.VaultsScreen()
			//case tcell.KeyDEL: //TODO удаление секрета
			//case tcell.Key: //TODO удаление секрета
		}
		return event
	})

	secretsTable.
		SetCell(0, 0, tview.NewTableCell("ID").SetSelectable(false)).
		SetCell(0, 1, tview.NewTableCell("Name").SetSelectable(false)).
		SetCell(0, 2, tview.NewTableCell("Type").SetSelectable(false)).
		SetCell(0, 3, tview.NewTableCell("Data").SetSelectable(false))

	secrets, err := t.provider.GetSecrets(vaultID)
	if err != nil {
		//TODO модальное окно с ошибками
	}
	for k, secret := range secrets {
		cellName := tview.NewTableCell(secret.Name).
			SetClickedFunc(func() bool {

				return true
			})
		cellType := tview.NewTableCell(secret.Type)
		cellData := tview.NewTableCell(fmt.Sprintf("[%s][%s]\n[%s]",
			secret.VaultID,
		)).SetExpansion(1)
		secretsTable.SetCellSimple(k+1, 0, secret.ID)
		secretsTable.SetCell(k+1, 1, cellName)
		secretsTable.SetCell(k+1, 2, cellType)
		secretsTable.SetCell(k+1, 3, cellData)
	}

	//cellComment := tview.NewTableCell(fmt.Sprintf("[%s]\n[%s]")).SetExpansion(1)
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(newPrimitive("Secrets"), 0, 10, false)
	flex.AddItem(secretsTable, 0, 100, true)
	flex.AddItem(newPrimitive("Esc for return to vaults, 'e' for edit, 'enter' for view secrets, 'delete' for delete secret"), 0, 10, false)

	t.App.SetRoot(flex, true)
}

func (t *TUI) VaultEditScreen() {}

func (t *TUI) SecretViewScreen(secretID string) {
	//todo: получать секрет по ID
	//todo: персональный вью для каждого типа секретов
	//todo: сделать таблицей
	secret, err := t.provider.GetSecret(secretID)
	if err != nil {
		//todo модальное окно с ошибкой
		return
	}
	secretTable := tview.NewTable()
	secretTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			t.SecretsScreen(secret.VaultID)
			//case tcell.KeyDEL: //TODO удаление секрета
			//case tcell.Key: //TODO Редактирование секрета
		}
		return event
	})
	r := 0
	secretTable.SetCellSimple(r, 0, Name).SetCellSimple(r, 1, secret.Name)
	r++
	secretTable.SetCellSimple(r, 0, Text).SetCellSimple(r, 1, secret.UnPackedSecret)
	r++
	secretTable.SetCellSimple(r, 0, ID).SetCellSimple(r, 1, secret.ID)
	r++
	secretTable.SetCellSimple(r, 0, VaultID).SetCellSimple(r, 1, secret.VaultID)
	r++
	secretTable.SetCellSimple(r, 0, Data).SetCellSimple(r, 1, string(secret.Secret))
	r++
	secretTable.SetCellSimple(r, 0, CreateAt).SetCellSimple(r, 1, secret.CreateAt.String())
	r++
	secretTable.SetCellSimple(r, 0, UpdateAt).SetCellSimple(r, 1, secret.UpdateAt.String())
	r++
	secretTable.SetCellSimple(r, 0, Comment).SetCellSimple(r, 1, secret.Comment)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(newPrimitive("Secret"), 0, 10, false)
	flex.AddItem(secretTable, 0, 100, true)
	flex.AddItem(newPrimitive("Esc for return to vaults, 'e' for edit, 'tab' for delete secretTable"), 0, 10, false)

	t.App.SetRoot(flex, true)
}
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
			case Secrets:
				t.SecretsScreen("")
			case VaultEdit:
				t.VaultEditScreen()
			case SecretView:
				t.SecretViewScreen("")
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
