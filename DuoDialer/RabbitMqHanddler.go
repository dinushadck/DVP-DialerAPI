package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"time"
)

func RabbitMQPublish(queueName string, publishData []byte) {
	// Connects opens an AMQP connection from the credentials in the URL.
	dialString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rabbitMQHost, rabbitMQPort)
	conn, err := amqp.Dial(dialString)
	if err != nil {
		fmt.Println("connection.open: ", err)
	}

	// This waits for a server acknowledgment which means the sockets will have
	// flushed all outbound publishings prior to returning.  It's important to
	// block on Close to not lose any publishings.
	defer conn.Close()

	c, err := conn.Channel()
	if err != nil {
		fmt.Println("channel.open: ", err)
	}

	// We declare our topology on both the publisher and consumer to ensure they
	// are the same.  This is part of AMQP being a programmable messaging model.
	//
	// See the Channel.Consume example for the complimentary declare.
	//err = c.ExchangeDeclare(queueName, "topic", true, false, false, false, nil)
	//if err != nil {
	//	fmt.Println("exchange.declare: ", err)
	//}

	// Prepare this message to be persistent.  Your publishing requirements may
	// be different.
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		ContentType:  "application/json",
		Body:         publishData,
	}

	// This is not a mandatory delivery, so it will be dropped if there are no
	// queues bound to the logs exchange.
	err = c.Publish("", queueName, false, false, msg)
	if err != nil {
		// Since publish is asynchronous this can happen if the network connection
		// is reset or if the server has run out of resources.
		fmt.Println("basic.publish: ", err)
	} else {
		fmt.Println("basic.publish: ", queueName)
	}
}
