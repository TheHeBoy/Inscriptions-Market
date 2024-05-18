package mqI

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"gohub/pkg/rabbitmq"
)

var channel *amqp.Channel

func SetupMQ() {
	rabbitmq.ConnectMQ()
	channel = rabbitmq.Ch
	List.Declare()
}
