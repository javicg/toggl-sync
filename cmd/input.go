package cmd

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

type InputController interface {
	RequestTextInput(string) (string, error)
	RequestPassword(string) (string, error)
}

type StdInController struct{}

func (StdInController) RequestTextInput(description string) (string, error) {
	fmt.Print(description)
	r := bufio.NewReader(os.Stdin)
	return r.ReadString('\n')
}

func (StdInController) RequestPassword(description string) (string, error) {
	fmt.Print(description)
	bytes, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	return string(bytes), err
}
