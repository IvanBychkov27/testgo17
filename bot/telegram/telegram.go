// https://www.youtube.com/watch?v=gegtZMx7poE

// https://api.telegram.org/bot5118317614:AAF3V4ry07OwYp7noV6PpG8d8Ptv3r3mUG0/getMe
// {"ok":true,"result":{"id":5118317614,"is_bot":true,"first_name":"IvBot","username":"SimpleDevelopmentIvBot","can_join_groups":true,"can_read_all_group_messages":false,"supports_inline_queries":false}}

// https://core.telegram.org/bots/api - документация к боту
package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

func main() {
	log.Println("Start IvBot...")

	botToken := "5118317614:AAF3V4ry07OwYp7noV6PpG8d8Ptv3r3mUG0"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken

	offSet := 0
	for {
		updates, err := getUpdate(botUrl, offSet)
		if err != nil {
			log.Println("error:", err.Error())
		}

		for _, update := range updates {
			err = respond(botUrl, update)
			if err != nil {
				log.Println("error:", err.Error())
			}
			offSet = update.UpdateID + 1
		}
		//fmt.Println(updates)
	}
}

// запрос обновлений
func getUpdate(botUrl string, offSet int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offSet))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse

	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}

	return restResponse.Result, nil
}

// ответ на обновления
func respond(botUrl string, update Update) error {
	var botMessage BotMessage

	botMessage.ChatID = update.Message.Chat.ChatID
	botMessage.Text = update.Message.Text

	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}

	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}

	return nil
}
