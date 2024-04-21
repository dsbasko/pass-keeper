package tui

const (
	Exit         = "exit"
	Logout       = "logout"
	StartMenu    = "startMenu"
	LoginMenu    = "loginScreen"
	RegisterMenu = "registerScreen"
	CreatePIN    = "createPIN"
	EnterPIN     = "enterPIN"
	MainMenu     = "mainMenu"
	Vaults       = "vaults"
	Vault        = "vault"
	VaultEdit    = "vaultEdit"
	SecretView   = "secret"
	SecretEdit   = "secretEdit"
)

var StartMenuTransitions = []string{Exit, RegisterMenu, LoginMenu}
var RegisterMenuTransition = []string{Exit, CreatePIN, StartMenu}
var LoginMenuTransition = []string{Exit, EnterPIN, StartMenu}
var CreatePINTransitions = []string{Exit, LoginMenu}
var EnterPINTransitions = []string{Exit, Logout, MainMenu}
var MainMenuTransitions = []string{Exit, Logout, Vaults}
var VaultsTransitions = []string{Exit, Logout, Vault, VaultEdit, MainMenu}
var VaultTransitions = []string{Exit, Vaults, VaultEdit}
var VaultEditTransitions = []string{Exit, LoginMenu, Vault}
var SecretViewTransitions = []string{Exit, Vault, SecretEdit}
var SecretEditTransitions = []string{Exit, Vault, SecretView}
