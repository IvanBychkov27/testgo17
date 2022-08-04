// https://dev-gang.ru/article/sozdanie-avtomaticzeskogo-raspisanija-planirovsczik-zadanii-cron-s-pomosczu-golang-lgc6kt5exf/?ysclid=l1ghntf8ec
// Создание автоматического расписания (планировщик заданий Cron) с помощью Golang
// Cron - это команда Linux для запуска скриптов по расписанию.
/*

* * * * * perintah yang akan dieksekusi
– – – – –
| | | | |
| | | | +—– day of week (0 – 7) (Sunday=0)
| | | +——- month (1 – 12)
| | +——— day of month (1 – 31)
| +———– hour (0 – 23)
+————- min (0 – 59)

*/

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	cronPkg "github.com/robfig/cron/v3"
)

func main() {
	//cron := cronPkg.New() // по минутный
	cron := cronPkg.New(cronPkg.WithSeconds()) // по секундный
	defer cron.Stop()

	ch := make(chan struct{}, 1)

	cronKey := "*/2 * * * * *"
	//cron.AddFunc(cronKey, func() { run("Ivan") })
	cron.AddFunc(cronKey, func() { ticker(ch) })

	go cron.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case <-ch:
			log.Println("channel ok")
		case <-sig:
			return
		}
	}

}

func run(d string) {
	fmt.Println(time.Now().Format("2006-01-02 15:04:05") + " Hello, " + d + "! I am Cron...")
}

func ticker(sig chan struct{}) {
	//log.Println("channel...")
	sig <- struct{}{}
}

// сбрасывает у списка юзеров в БД флаг online
//func (app *Application) ClearUser(ctx context.Context, wg *sync.WaitGroup) {
//	defer wg.Done()
//
//	cron := cronPkg.New()
//	_, errCron := cron.AddFunc(app.cfg.ClearUsers, app.clear)
//	if errCron != nil {
//		app.logger.Error("error add func to cron", zap.Error(errCron))
//		return
//	}
//
//	cron.Start()
//
//	<-ctx.Done()
//
//	cron.Stop()
//}
