package tui

import (
	"fmt"
	"strings"
)

func (c *Commands) Help() error {
	var out strings.Builder

	out.WriteString("Commands:")
	out.WriteString("\n  sign-up - register new user")
	out.WriteString("\n  sign-in - login user")
	out.WriteString("\n  -------")
	out.WriteString("\n  clear - clear screen")
	out.WriteString("\n  exit - exit from application")
	out.WriteString("\n  help - show help")

	fmt.Println(out.String())
	return nil
}
