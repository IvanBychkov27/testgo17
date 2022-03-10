package main

import (
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

// Queue config
type Queue struct {
	Key      string `env:"KEY"`
	Secret   string `env:"SECRET"`
	Endpoint string `env:"ENDPOINT"`
	URL      string `env:"URL"`
	Region   string `env:"REGION"`
}

func queueParams() Queue {
	return Queue{
		Key:      "adIPS2V8TcN8W3EierJW",
		Secret:   "f_F21HW-adAsYoAxAhecxHMdViLU2xjC4_DOKHoH",
		Endpoint: "https://message-queue.api.cloud.yandex.net",
		//URL:      "https://message-queue.api.cloud.yandex.net/b1g8fp34kitmv2q6kesr/dj6000000005cpe706ht/test2", // test2
		URL: "https://message-queue.api.cloud.yandex.net/b1g8fp34kitmv2q6kesr/dj6000000001uq3i06ht/test", // test
		//URL:    "https://message-queue.api.cloud.yandex.net/b1g8fp34kitmv2q6kesr/dj6000000005a44r06ht/offer-check", // боевой offer-check
		Region: "ru-central1",
	}
}

func main() {

	sendMessages() // отправить сообщения в очередь

	//del := true // удалить
	//del := false        // не удалять
	//receiveMessage(del) // если истина, то прочитанное сообщение будет удалено из очереди
}

// послать сообщения
func sendMessages() {
	msgs := []string{
		//`{"url":"http://4000.99.adscompass.ru","country_list": ["DK", "CA", "FR", "IT", "BE", "SE"],"pattern":"","platform_list":["Windows"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
		//`{"url":"http://4000.99.adscompass.ru","country_list": ["CA", "DK", "FR"],"pattern":"testKey", "platform_list":["Windows","Linux"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
		`{"url":"http://4000.99.adscompass.ru","country_list": ["CA"],"pattern":"testKey", "platform_list":["Windows"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
		//`{"url":"http://3000.99.adscompass.ru","country_list": ["US", "PL", "JP", "AT", "GB", "BR"],"pattern":"testKey","platform_list":["Linux"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
		//`{"url":"http://3000.99.adscompass.ru","country_list": ["CA", "DK"],"pattern":"testKey","platform_list":["Linux"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
		//`{"url":"https://nataliedate.com/wizard-man?utm_source=GoldLead&linkid=42224&clickid={clickid}&web_id={web_id}&sub_id={sub_id}","country_list": ["RU"],"pattern":"testKey","platform_list":["Linux"],"tags": {"offer_id": 6, "url_type": "main", "url_name": "link4"}}`,
	}
	for _, msg := range msgs {
		sendMessage(msg)
	}
}

// послать сообщение
func sendMessage(msg string) {
	q := queueParams()
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(q.Key, q.Secret, ""),
		Endpoint:    aws.String(q.Endpoint),
		Region:      aws.String(q.Region),
	})
	if err != nil {
		log.Printf("error create queue, %s", err.Error())
		return
	}

	queue := sqs.New(sess)

	input := &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(q.URL),
	}

	_, err = queue.SendMessage(input)
	if err != nil {
		log.Printf("error send messages to queue, %s", err.Error())
		return
	}

	log.Println("сообщение отправлено в очередь")
}

// получить сообщение
func receiveMessage(del bool) {
	q := queueParams()
	sess, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(q.Key, q.Secret, ""),
		Endpoint:    aws.String(q.Endpoint),
		Region:      aws.String(q.Region),
	})
	if err != nil {
		log.Printf("error create queue, %s", err.Error())
		return
	}

	queue := sqs.New(sess)

	input := &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(q.URL),
		MaxNumberOfMessages: aws.Int64(1), // кол-во возвращаемых сообщений от 1 до 10
		VisibilityTimeout:   aws.Int64(1), // время (в сек) блокировки сообщения - в течение этого времени это сообщение будет не доступно другим клиентам
		WaitTimeSeconds:     aws.Int64(1), // время (в сек) ожидания сообщений, если в течение этого времени сообщения нет, то отдает пустое сообщение
	}
	res, err := queue.ReceiveMessage(input)
	if err != nil {
		log.Printf("error get message from queue, %s", err.Error())
		return
	}

	out := res.String()
	if len(out) == 4 {
		log.Println("message nil")
		return
	}
	//log.Println("message:", out)

	log.Println("читаем из очереди...")
	for i, m := range res.Messages {
		log.Printf("сообщение №%d: %s \n", i+1, *m.Body)
		if del {
			delMessage(queue, q.URL, *m.ReceiptHandle)
		}
	}

}

// удаляет сообщение из очереди
func delMessage(queue *sqs.SQS, queueUrl string, receiptHandle string) {
	input := &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueUrl),
		ReceiptHandle: aws.String(receiptHandle),
	}

	_, err := queue.DeleteMessage(input)
	if err != nil {
		log.Printf("error del message from queue, %s", err.Error())
		return
	}

	log.Println("сообщение удалено из очереди")
}
