package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("MyDarkSecret")

	// Хэширование пароля со стоимостью по умолчанию 10
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(hashedPassword))

	// Сравнение пароля с хэшем
	err = bcrypt.CompareHashAndPassword(hashedPassword, password)
	fmt.Println(err) // nil означает, что это совпадение
}
