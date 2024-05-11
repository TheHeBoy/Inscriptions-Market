package mqI

import (
	"encoding/json"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"gohub/internal/model"
	"gohub/internal/service"
	"gohub/pkg/logger"
)

type ListService struct {
}

var List = new(ListService)

func (*ListService) Declare() {
	var exchangeName = "logs_topic"
	var routingKey = "list"
	var queueName = "listQueue"
	err := channel.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		panic(err)
	}

	q, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		panic(err)
	}

	err = channel.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil)
	if err != nil {
		panic(err)
	}

	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // args
	)
	if err != nil {
		panic(err)
	}
	consume(msgs)
}

func consume(msgs <-chan amqp.Delivery) {
	go func() {
		for d := range msgs {
			logger.Debugf("Received a List message:%s", d.Body)
			var list model.ListDO
			err := json.Unmarshal(d.Body, &list)
			if err != nil {
				logger.Errorf("Failed to unmarshal: %+v", errors.WithStack(err))
				continue
			}
			service.Order.Listing(list)
		}
	}()
}
