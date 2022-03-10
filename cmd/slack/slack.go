// https://onedev.net/post/935
/*
Порядок действий.

Открываем страницу приложения Incoming WebHooks - http://slack.com/apps/A0F7XDUAZ.
Открываем страницу добавления новой конфигурации хука "Add Configuration".
Выбираем канал для уведомлений и жмем "Add Incoming WebHooks integration".
Получаем URL для отправки POST запроса с данными playload JSON.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var message string

	//hookUrl := "https://hooks.slack.com/services/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"
	hookUrl := "https://hooks.slack.com/services/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"
	message = "message bot iv"

	playload := map[string]interface{}{
		"username": "bot-iv",
		"text":     message,
	}

	bodyString, _ := sendPostJson(hookUrl, playload)
	log.Println("body", bodyString)
}

func sendPostJson(url string, data interface{}) (bodyString string, err error) {
	info, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(info))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	bodyString = string(body)

	if resp.StatusCode != 200 {
		fmt.Println("Response Status:", resp.Status) // 200 OK
		fmt.Println("Response Headers:", resp.Header)
		fmt.Println("Response Body:", bodyString)

		return bodyString, fmt.Errorf("error %w", err)
	}

	return bodyString, nil
}
