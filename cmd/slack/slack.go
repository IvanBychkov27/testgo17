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
	//hookUrl := "https://hooks.slack.com/services/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"
	hookUrl := "https://hooks.slack.com/services/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"
	message = "message bot iv"

	playload := map[string]interface{}{
		"username": "bot-iv5",
		"text":     message,
		//"icon_emoji": ":heavy_check_mark:",
		"icon_emoji": ":ghost:",
	}

	/* см список эмодзи:
	   https://slackmojis.com/
	   https://lesswrong.ru/wiki/%D0%A1%D0%BF%D0%B8%D1%81%D0%BE%D0%BA_%D1%8D%D0%BC%D0%BE%D0%B4%D0%B7%D0%B8_Slack-%D1%87%D0%B0%D1%82%D0%B0?ysclid=l6bymx604v630507163
	*/

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
