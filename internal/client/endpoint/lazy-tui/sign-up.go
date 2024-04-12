package tui

import (
	"fmt"
	"syscall"

	"golang.org/x/term"

	"github.com/dsbasko/pass-keeper/pkg/errors"
)

func (c *Commands) SignUp() func() error {
	return func() (err error) {
		defer errors.ErrorPtrWithOP(&err, "endpoint.lazy-tui.SignUp")
		var email, password, secretKey string

		fmt.Printf("Enter email:\n")
		if _, err = fmt.Scanln(&email); err != nil {
			return errors.ErrorWithOP(err, "failed to read email")
		}

		fmt.Printf("\nEnter password:\n")
		bytePassword, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return errors.ErrorWithOP(err, "failed to read password")
		}
		password = string(bytePassword)

		fmt.Printf("\nEnter secret-key:\n")
		if _, err = fmt.Scanln(&secretKey); err != nil {
			return errors.ErrorWithOP(err, "failed to read secret-key")
		}

		fmt.Printf("\nEmail: %s\n", email)
		fmt.Printf("Password: %s\n", password)
		fmt.Printf("Secret-key: %s\n", secretKey)

		return nil
	}
}
