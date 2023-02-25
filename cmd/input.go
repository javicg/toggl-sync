package cmd

import (
	"bufio"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

type inputController interface {
	requestTextInput(string) (string, error)
	requestPassword(string) (string, error)
}

// StdInController is an input controller that redirects all calls to Stdin
type StdInController struct{}

func (StdInController) requestTextInput(description string) (string, error) {
	fmt.Print(description)
	r := bufio.NewReader(os.Stdin)
	return r.ReadString('\n')
}

func (StdInController) requestPassword(description string) (string, error) {
	fmt.Print(description)
	bytes, err := term.ReadPassword(syscall.Stdin)
	fmt.Println()
	return string(bytes), err
}
