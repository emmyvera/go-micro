package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setUp()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil

}

func (Consumer *Consumer) setUp() error {
	channel, err := Consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declearExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (Consumer *Consumer) Listen(topics []string) error {
	ch, err := Consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := declearQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchage, Queue] [logs_topic, %s]\n", q.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Print(err)
		}

	case "auth":
		// authenticate

		// we can have as many cases as possible provided with the logic

	default:
		err := logEvent(payload)
		if err != nil {
			log.Print(err)
		}
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Print(err)
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Print(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		log.Print(err)
		return err
	}

	return nil
}
