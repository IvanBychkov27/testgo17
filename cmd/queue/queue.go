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

func main() {

	sendMessages() // отправить сообщения в очередь

	//del := true // удалить
	//del := false        // не удалять
	//receiveMessage(del) // если истина, то прочитанное сообщение будет удалено из очереди
}

// послать сообщения
func sendMessages() {
	msgs := []string{
		//`{"1":"http://4000.99.adscompass.ru","2": ["DK", "CA"]}`,
		//`{"1":"http://4000.99.adscompass.ru","2": ["CA"]}`,
		//`{"1":"http://4000.99.adscompass.ru","2": ["CA","PL"]}`,
		`{"1":"http://4000.99.adscompass.ru","2": ["DK"]}`,
	}
	for _, msg := range msgs {
		sendMessage(msg)
	}
}

// послать сообщение
func sendMessage(msg string) {
	q := Queue{
		Key:      "adIPS2V8TcN8W3EierJW",
		Secret:   "f_F21HW-adAsYoAxAhecxHMdViLU2xjC4_DOKHoH",
		Endpoint: "https://message-queue.api.cloud.yandex.net",
		URL:      "https://message-queue.api.cloud.yandex.net/b1g8fp34kitmv2q6kesr/dj6000000001uq3i06ht/test",
		Region:   "ru-central1",
	}

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
	q := Queue{
		Key:      "adIPS2V8TcN8W3EierJW",
		Secret:   "f_F21HW-adAsYoAxAhecxHMdViLU2xjC4_DOKHoH",
		Endpoint: "https://message-queue.api.cloud.yandex.net",
		URL:      "https://message-queue.api.cloud.yandex.net/b1g8fp34kitmv2q6kesr/dj6000000001uq3i06ht/test",
		Region:   "ru-central1",
	}
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
