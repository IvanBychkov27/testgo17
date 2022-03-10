package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func main() {
	username, password, _ := credentials()
	fmt.Println()
	fmt.Printf("Username: %s, Password: %s\n", username, password)
}

func credentials() (string, string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter Username: ")

	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}
	fmt.Print("Enter Password: ")

	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)

	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}
