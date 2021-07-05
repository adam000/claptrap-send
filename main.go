package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

const usage = `claptrap-send

Usage:
    claptrap-send --user=<user> --host=<host> --port=<port> --vhost=<vhost> --topic=<topic> --from=<from> --subject=<subject> --message=<message>
`

const passwordEnvVarName = "RABBITMQ_PASSWORD"

type Message struct {
	From    string
	Subject string
	Body    string
}

func main() {
	arguments, err := docopt.ParseDoc(usage)
	if err != nil {
		log.Fatalf("Failed to parse args: %v", err)
	}

	password := os.Getenv(passwordEnvVarName)
	if password == "" {
		log.Printf("No environment variable set at %s, attempting with blank password", passwordEnvVarName)
	}

	username, err := arguments.String("--user")
	if err != nil {
		log.Fatalf("Failed to get argument 'user': %v", err)
	}
	host := arguments["--host"]
	port := arguments["--port"]
	vhost := arguments["--vhost"]
	topic, err := arguments.String("--topic")
	if err != nil {
		log.Fatalf("Failed to get argument 'topic': %v", err)
	}
	from, err := arguments.String("--from")
	if err != nil {
		log.Fatalf("Failed to get argument 'from': %v", err)
	}
	subject, err := arguments.String("--subject")
	if err != nil {
		log.Fatalf("Failed to get argument 'subject': %v", err)
	}
	body, err := arguments.String("--message")
	if err != nil {
		log.Fatalf("Failed to get argument 'message': %v", err)
	}

	//log.Printf("%s %s %s %s %s %s", username, password, host, port, vhost, topic)
	//log.Printf("%s %s %s", from, subject, body)

	connString := fmt.Sprintf("amqp://%s:%s@%s:%s//%s", username, password, host, port, vhost)
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer channel.Close()

	/*
		queue, err := channel.QueueDeclare(
			topic, // name
			true,  // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   //arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare a queue: %v", err)
		}
	*/

	message := Message{
		From:    from,
		Subject: subject,
		Body:    body}

	bytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message: %v", err)
	}

	err = channel.Publish(
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        bytes,
		})

	if err != nil {
		log.Fatalf("Failed to publish the message: %v", err)
	}

	log.Println("Message sent")
}
