// https://medium.com/golang-notes/%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%B0-%D1%81-%D0%BE%D0%B4%D0%BD%D0%BE%D1%80%D0%B0%D0%B7%D0%BE%D0%B2%D1%8B%D0%BC%D0%B8-%D0%BF%D0%B0%D1%80%D0%BE%D0%BB%D1%8F%D0%BC%D0%B8-%D0%B2-go-5b1c5b10fc23

package main

import (
	"bytes"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"html/template"
	"image/png"
	"log"
	"net/http"
)

//тут будем хранить TOTP для одного пользователя
var key *otp.Key

func main() {
	//Настраиваем TOTP
	//для каждого пользователя TOTP ключ должен быть уникальным
	//В нашей программе ключ будет разный с каждым запуском (!)
	var err error
	key, err = totp.Generate(totp.GenerateOpts{
		Issuer:      "Example.com",
		AccountName: "Ivan@example.com",
	})
	if err != nil {
		panic(err)
	}

	log.Println("start addr :3000")
	http.HandleFunc("/", indexHandlerFunc)
	http.HandleFunc("/login/", loginHandlerFunc)
	http.HandleFunc("/2fa/", setup2FAHandlerFunc)
	http.HandleFunc("/qr.png", genQRCodeHandlerFunc)
	http.ListenAndServe(":3000", nil)
}

func indexHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//для простоты не обрабатываем ошибки
	t, _ := template.ParseFiles("cmd/password/web/templates/index.html")
	t.Execute(w, nil)
}

func loginHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//Обрабатываем только POST-запрос
	if r.Method != "POST" {
		http.NotFound(w, r)
	}
	//для простоты не обрабатываем ошибки
	r.ParseForm()
	user := r.FormValue("user")

	password := r.FormValue("password") //Проверяем логин и пароль
	if !(user == "Ivan" && password == "123") {
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		return
	}

	w.Write([]byte("Hello, " + user))
}

//Отображает страницу с QR-кодом
func setup2FAHandlerFunc(w http.ResponseWriter, r *http.Request) {
	//для простоты не обрабатываем ошибки
	t, _ := template.ParseFiles("cmd/password/web/templates/2fa.html")
	t.Execute(w, nil)
}

//Генерирует QR-код для добавления аккаунта в Яндекс.Ключ/Google.Authentificator
func genQRCodeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	//для простоты не обрабатываем ошибки
	png.Encode(&buf, img)
	w.Header().Set("Content-Type", "image/png")
	w.Write(buf.Bytes())
}
