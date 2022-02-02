package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Println(time.Now())
	now := time.Now().UTC().Unix()

	fmt.Println("now    =", now)

	d := time.Unix(now, 0)
	fmt.Println("data =", d)

	fmt.Println("===========================")

	//var timeStart, timeStop  time.Time
	timeStart := time.Date(2021, 12, 10, 7, 0, 0, 0, time.UTC)
	timeStop := time.Date(2021, 12, 14, 7, 0, 0, 0, time.UTC)
	fmt.Println("timeStart =", timeStart)
	fmt.Println("timeStop  =", timeStop)

	dif := timeStop.Sub(timeStart)
	fmt.Println("dif =", dif)
	//newTime := timeStop.Add(dif)
	newTime := timeStop.Add(-time.Hour * 24 * 4)
	fmt.Println("newTime  =", newTime)

	fmt.Println()
	fmt.Println("===========================")

	dateList := setDateList(timeStart, timeStop)
	fmt.Println("dateList  =", dateList)

}

func setDateList(from, to time.Time) map[string]time.Time {
	dateList := make(map[string]time.Time)
	dateList[from.Format("2006-01-02")] = from
	dateList[to.Format("2006-01-02")] = to

	dif := to.Unix() - from.Unix()
	if dif <= 0 {
		return dateList
	}

	for {
		from = from.Add(time.Hour * 24)
		key := from.Format("2006-01-02")
		_, ok := dateList[key]
		if ok {
			break
		}
		dateList[key] = from
	}

	return dateList
}

// timeNow функция выводит текущую дату и время, addHour позволяет прибавить/отнять часы (для корректировки)
func timeNow(addHour int) string {
	y := time.Now().Year()
	mec := time.Now().Month()
	d := time.Now().Day()
	h := time.Now().Hour() + addHour
	m := time.Now().Minute()
	s := time.Now().Second()
	return time.Date(y, mec, d, h, m, s, 0, time.UTC).Format("02.01.2006  15:04:05")
}
