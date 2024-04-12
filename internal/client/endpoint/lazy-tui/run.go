package tui

import (
	"context"
	"fmt"
	"os"
)

type Commands struct{}

func Run(ctx context.Context) error {
	c := &Commands{}
	c.greeting()

	for {
		inputCh, errCh := c.Scan()
		select {
		case <-ctx.Done():
			return nil
		case <-errCh:
			return nil
		case input := <-inputCh:
			c.runCommand(input)
		}
	}
}

func (c *Commands) greeting() {
	fmt.Println("Welcome to PassKeeper!")
	fmt.Println("Please sign-up or sign-in to continue...")
	fmt.Println("Type 'help' to see available commands...")
	fmt.Println()
}

func (c *Commands) runCommand(input string) {
	switch input {
	case "signup", "sign-up", "register", "reg":
		c.Wrap(c.SignUp())
	case "signin", "sign-in", "login":
		c.Wrap(c.SignIn)
	case "clear":
		c.clear()
		c.greeting()
	case "exit", "quit", "q":
		os.Exit(0)
	case "help", "h":
		c.Wrap(c.Help)
	default:
		fmt.Println()
		fmt.Println("Unknown command")
		c.Wrap(c.Help)
	}
}

func (c *Commands) Wrap(fn func() error) {
	fmt.Println()
	if err := fn(); err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println()
}

func (c *Commands) Scan() (chan string, chan error) {
	chInput := make(chan string)
	chError := make(chan error)

	go func() {
		defer close(chInput)
		defer close(chError)

		var input string
		if _, err := fmt.Scanln(&input); err != nil {
			chError <- fmt.Errorf("scan error: %w", err)
			return
		}

		chInput <- input
	}()

	return chInput, chError
}
