package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gabrielmellooliveira/order-service/pkg/events"
	"github.com/streadway/amqp"
	"sync"
)

type OrderCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewOrderCreatedHandler(rabbitMQChannel *amqp.Channel) *OrderCreatedHandler {
	return &OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *OrderCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) error {
	defer wg.Done()

	fmt.Printf("Order created: %v\n", event.GetPayload())

	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		return err
	}

	msgRabbitMQ := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	err = h.RabbitMQChannel.Publish(
		"amq.direct",
		"order_created",
		false,
		false,
		msgRabbitMQ,
	)
	if err != nil {
		return err
	}

	return nil
}
