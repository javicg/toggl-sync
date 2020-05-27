package cmd

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
)

type inputController interface {
	requestTextInput(string) (string, error)
	requestPassword(string) (string, error)
}

type stdInController struct{}

func (stdInController) requestTextInput(description string) (string, error) {
	fmt.Print(description)
	r := bufio.NewReader(os.Stdin)
	return r.ReadString('\n')
}

func (stdInController) requestPassword(description string) (string, error) {
	fmt.Print(description)
	bytes, err := terminal.ReadPassword(syscall.Stdin)
	fmt.Println()
	return string(bytes), err
}
