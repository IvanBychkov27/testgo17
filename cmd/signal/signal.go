// https://medium.com/nuances-of-programming/%D0%BE%D0%B1%D1%80%D0%B0%D0%B1%D0%BE%D1%82%D0%BA%D0%B0-%D1%81%D0%B8%D0%B3%D0%BD%D0%B0%D0%BB%D0%BE%D0%B2-%D0%B2-%D0%BE%D0%BF%D0%B5%D1%80%D0%B0%D1%86%D0%B8%D0%BE%D0%BD%D0%BD%D1%8B%D1%85-%D1%81%D0%B8%D1%81%D1%82%D0%B5%D0%BC%D0%B0%D1%85-%D1%81%D0%B5%D0%BC%D0%B5%D0%B9%D1%81%D1%82%D0%B2%D0%B0-unix-%D0%BD%D0%B0-golang-cb2c42b80ba5
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("start...")

	// SIGTERM — это общий сигнал, используемый для завершения программы
	// Сигнал SIGINT отправляется при введении пользователем в управляющем терминале символа прерывания, по умолчанию это ^C (Control-C)
	// Сигнал SIGQUIT отправляется при введении пользователем в управляющем терминале символа выхода, по умолчанию это ^\ (Control-Backslash)
	// Сигнал SIGHUP отправляется при потере программой своего управляющего терминала

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	dir, err := os.UserHomeDir()
	if err != nil {
		log.Println("error user home dir", err.Error())
	}
	log.Println("user home dir:", dir)

	temp := os.TempDir()
	log.Println("temp:", temp)

	c, err := os.UserCacheDir()
	if err != nil {
		log.Println("error user cache dir", err.Error())
	}
	log.Println("user cache dir:", c)

	conf, err := os.UserConfigDir()
	if err != nil {
		log.Println("error user config dir", err.Error())
	}
	log.Println("user config dir:", conf)

	<-ctx.Done()
	log.Println("done")
}
