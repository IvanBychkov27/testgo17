package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

// redis_read - только для чтения данных
func main() {
	// подключение к Redis
	redisDB, err := redis.Dial("tcp", "127.0.0.1:6379") // данный порт устанавливается по умолчанию для Redis при его установке
	if err != nil {
		fmt.Println("connect to redis error", err)
		return
	} else {
		fmt.Println("connection to redis...")
	}
	defer redisDB.Close()

	password := "pass"
	_, err = redisDB.Do("AUTH", password) // ввод пароля для redis
	if err != nil {
		fmt.Println("error password to redis", err.Error())
		return
	}

	readData(redisDB)

	fmt.Println("Done redis_read...")
}

func readData(redisDB redis.Conn) {
	exit := "c"
	fmt.Println("для выхода наберите -", exit)
	var s string
	for {
		if s == "" {
			s = "userbalances"
		}
		res, err := readRedis(redisDB, s)
		if err == nil {
			if len(res) > 500 {
				fmt.Printf("%s:\n%s \n", s, res[:500])
			} else {
				fmt.Printf("%s:\n%s \n", s, res)
			}
		} else {
			if err.Error() == "redigo: nil returned" {
				fmt.Println("redis get failed: nil")
			} else {
				fmt.Println("error redis get failed:", err.Error())
			}
		}

		fmt.Scanln(&s)
		if s == exit {
			return
		}
	}
}

// чтение данных
func readRedis(redisDB redis.Conn, key string) (res string, err error) {
	res, err = redis.String(redisDB.Do("GET", key))
	if err != nil {
		return "", err
	}
	return res, nil
}
