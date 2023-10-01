package usecase

import (
	"github.com/gabrielmellooliveira/order-service/internal/entity"
	"github.com/gabrielmellooliveira/order-service/pkg/events"
)

type OrderInputDto struct {
	ID    string  `json:"id"`
	Price float64 `json:"price"`
	Tax   float64 `json:"tax"`
}

type OrderOutputDto struct {
	ID         string  `json:"id"`
	Price      float64 `json:"price"`
	Tax        float64 `json:"tax"`
	FinalPrice float64 `json:"finalPrice"`
}

type CreateOrderUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewCreateOrderUseCase(
	orderRepository entity.OrderRepositoryInterface,
	orderCreated events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (u *CreateOrderUseCase) Execute(input OrderInputDto) (OrderOutputDto, error) {
	order, err := entity.NewOrder(input.ID, input.Price, input.Tax)
	if err != nil {
		return OrderOutputDto{}, err
	}

	err = order.CalculateFinalPrice()
	if err != nil {
		return OrderOutputDto{}, err
	}

	err = u.OrderRepository.Save(order)
	if err != nil {
		return OrderOutputDto{}, err
	}

	output := OrderOutputDto{
		ID:         order.ID,
		Price:      order.Price,
		Tax:        order.Tax,
		FinalPrice: order.FinalPrice,
	}

	u.OrderCreated.SetPayload(output)
	u.EventDispatcher.Dispatch(u.OrderCreated)

	return output, nil
}
