package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/streadway/amqp"
)

func RabbitMQPublish(queueName string, publishData []byte) {
	if useAmqpAdapter == "true" {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in publish to amqp adapter", r)
			}
		}()

		//upload to amqp adapter service
		serviceurl := fmt.Sprintf("http://localhost:%s/DVP/API/1.0.0.0/amqpPublish", amqpAdapterPort)

		req, err := http.NewRequest("POST", serviceurl, bytes.NewBuffer(publishData))
		req.Header.Set("Content-Type", "application/json")
		fmt.Println("request:", serviceurl)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer resp.Body.Close()

		body, errb := ioutil.ReadAll(resp.Body)
		if errb != nil {
			fmt.Println(err.Error())
		} else {
			result := string(body)
			fmt.Println("response Body:", result)
		}
	} else {
		// Connects opens an AMQP connection from the credentials in the URL.
		rmqIps := strings.Split(rabbitMQHost, ",")
		defaultAmqpIP := rmqIps[0]

		dialString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, defaultAmqpIP, rabbitMQPort)
		conn, err := amqp.Dial(dialString)
		if err != nil {
			fmt.Println("connection1.open: ", err)
			if len(rmqIps) > 1 {
				dialString = fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rmqIps[1], rabbitMQPort)
				conn, err = amqp.Dial(dialString)

				if err != nil {
					fmt.Println("connection2.open: ", err)
				}
			}
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
		}
	}
}


func DashboardRabbitMQPublish(queueName string, publishData []byte) {

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in publish to amqp adapter", r)
		}
	}()

	rmqIps := strings.Split(rabbitMQHost, ",")
		defaultAmqpIP := rmqIps[0]

		dialString := fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, defaultAmqpIP, rabbitMQPort)
		conn, err := amqp.Dial(dialString)
		if err != nil {
			fmt.Println("connection1.open: ", err)
			if len(rmqIps) > 1 {
				dialString = fmt.Sprintf("amqp://%s:%s@%s:%s/", rabbitMQUser, rabbitMQPassword, rmqIps[1], rabbitMQPort)
				conn, err = amqp.Dial(dialString)

				if err != nil {
					fmt.Println("connection2.open: ", err)
				}
			}
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
			fmt.Println("dashboard.publish: ", err)
		} else {
			fmt.Println("dashboard.publish: ", queueName)
		}
}
