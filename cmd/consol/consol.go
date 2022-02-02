package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	err := consol()
	if err != nil {
		fmt.Println("error:", err.Error())
	} else {
		fmt.Println("ok")
	}

	fmt.Println("Done...")
}

func consol() error {
	// Go to the repo's directory
	r := "cmd/consol"
	err := os.Chdir(r)
	if err != nil {
		return err
	}
	// Print the command
	fmt.Printf("[%s] \n", r)

	// Execute the command
	out, errCMD := exec.Command("ls").CombinedOutput()
	//out, errCMD := exec.Command("pwd").CombinedOutput()
	//out, errCMD := exec.Command("mkdir", "myDir").CombinedOutput() // создать коталог
	//out, errCMD := exec.Command("rmdir", "myDir").CombinedOutput() // удалить каталог
	if errCMD != nil {
		fmt.Println("error command:", errCMD.Error())
	}

	fmt.Println(string(out))
	return nil
}
