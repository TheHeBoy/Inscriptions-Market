package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gohub/pkg/config"
)

var conn *amqp.Connection
var Ch *amqp.Channel

func ConnectMQ() {
	mqUrl := config.Get("mq.url")
	var err error
	conn, err = amqp.Dial(mqUrl)
	if err != nil {
		panic(err)
	}

	Ch, err = conn.Channel()
	if err != nil {
		panic(err)
	}

}
