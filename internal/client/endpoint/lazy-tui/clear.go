package tui

import (
	OS "os"
	"os/exec"
	"runtime"
)

func (c *Commands) clear() {
	switch os := runtime.GOOS; os {
	case "linux", "darwin":
		cmd := exec.Command("clear")
		cmd.Stdout = OS.Stdout
		_ = cmd.Run()
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = OS.Stdout
		_ = cmd.Run()
	default:
		print("\033[H\033[2J")
	}
}
