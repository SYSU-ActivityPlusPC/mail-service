package service

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"github.com/sysu-activitypluspc/mail-service/types"
)

// HandleMessage reveive msg from mq and let it mail
func HandleMessage() {
	forever := make(chan bool)

	// Get mq detaild message
	addr := os.Getenv("MQ_ADDRESS")
	if len(addr) == 0 {
		addr = "localhost"
	}
	port := os.Getenv("MQ_PORT")
	if len(port) == 0 {
		port = "5672"
	}
	user := os.Getenv("MQ_USER")
	pass := os.Getenv("MQ_PASSWORD")
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, addr, port)

	// Connect to mq
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Finished connecting to the mq")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ch.Qos(
		5,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		for d := range msgs {
			emailContent := new(types.EmailContent)
			err := json.Unmarshal(d.Body, emailContent)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = SendMail(emailContent.From, emailContent.To, emailContent.Content, emailContent.Subject)
			if err != nil {
				fmt.Println(err)
				continue
			}
			d.Ack(false)
		}
	}()

	<-forever
}
