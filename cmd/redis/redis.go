// https://russianblogs.com/article/84191656529/
// Использование и работа Go Language Redis
// Redis - это открытый исходный код, написанный на языке C, поддерживающий сетевое взаимодействие, постоянную базу данных Key-Value на основе памяти.

// https://timeweb.com/ru/community/articles/ustanovka-i-nastroyka-redis-1  - Установка и настройка Redis
// https://www.dmosk.ru/miniinstruktions.php?mini=redis-ubuntu  - Установка и настройка Redis + Docker

package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
)

func main() {
	// подключение к Redis
	redisDB, err := redis.Dial("tcp", "127.0.0.1:6379") // данный порт устанавливается по умолчанию для Redis при его установке
	if err != nil {
		fmt.Println("error connect to redis", err.Error())
		return
	}
	defer redisDB.Close()

	password := "pass"
	_, err = redisDB.Do("AUTH", password) // ввод пароля для redis
	if err != nil {
		fmt.Println("error password to redis", err.Error())
		return
	}

	fmt.Println("connection to redis...")

	writeRedis(redisDB, "test")

	//writeBatch(redisDB)

	readData(redisDB)

	fmt.Println("Done...")
}

func readData(redisDB redis.Conn) {
	exit := "c"
	fmt.Println("для выхода наберите -", exit)
	var s string
	for {
		res, err := readRedis(redisDB)
		if err == nil {
			fmt.Println("user:", res)
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
func readRedis(redisDB redis.Conn) (res string, err error) {
	res, err = redis.String(redisDB.Do("GET", "user"))
	if err != nil {
		return "", err
	}
	return res, nil
}

// запись данных
func writeRedis(redisDB redis.Conn, data string) {
	_, err := redisDB.Do("SET", "user", data, "EX", "10") // "EX", "120" - это значит что запись самоудалится через 120 сек (можно не указывать, тогда запись будет оставаться в памяти)
	if err != nil {
		fmt.Println("error set failed:", err.Error())
	}
}

// пакетная запись данных - ключем является весь слайс!
func writeBatch(redisDB redis.Conn) {
	key := []string{"a", "b", "c", "d", "e"}
	val := []string{"1", "2", "3", "4", "5"}
	_, err := redisDB.Do("MSET", key, val)
	if err != nil {
		fmt.Println("redis set failed:", err.Error())
	}
	v, e := redis.ByteSlices(redisDB.Do("MGET", key))
	if e != nil {
		fmt.Println("error redis set failed:", e.Error())
	} else {
		fmt.Println(string(v[0]))
	}

	//v, e :=redis.StringMap(redisDB.Do("MGET", key))
	//if e != nil {
	//	fmt.Println("error redis set failed:", e.Error())
	//} else {
	//	fmt.Println(v)
	//}

}

// Проверьте, существует ли значение ключа
// is_key_exit, err := redisDB.Bool(c.Do("EXISTS", "mykey1"))
// if err != nil {
// fmt.Println("error:", err)
// } else {
// fmt.Printf("exists or not: %v \n", is_key_exit)
// }

// Удалить ключ
// _, err = redisDB.Do("DEL", "mykey")
// if err != nil {
// fmt.Println("redis delelte failed:", err)
// }

// Установить время истечения
// n, _ := redisDB.Do("EXPIRE", key, 24*3600)
// if n == int64(1) {
// fmt.Println("success")
// }

/*
Метод Do объединяет методы Send, Flush и Receive.
Метод Do сначала записывает команду, затем очищает буфер вывода и, наконец,
получает все ожидающие ответы, включая результат команды, выданной стороной Do.
Если какой-либо ответ содержит ошибку, Do возвращает ошибку. Если ошибок нет, метод Do возвращает последний ответ.

c.Send ("SET", "foo", "bar") // Отправить записывает команды в подключенный выходной буфер
c.Send("GET", "foo")//
c.Flush () // Flush очищает подключенный выходной буфер и записывает его на сервер
// Recevie читает ответ сервера в порядке FIFO
c.Receive() // reply from SET
v, err = c.Receive() // reply from GET
*/
