package usecase

import (
	"github.com/gabrielmellooliveira/order-service/internal/entity"
)

type ListOrdersUseCase struct {
	OrderRepository entity.OrderRepositoryInterface
}

func NewListOrdersUseCase(
	orderRepository entity.OrderRepositoryInterface,
) *ListOrdersUseCase {
	return &ListOrdersUseCase{
		OrderRepository: orderRepository,
	}
}

func (u *ListOrdersUseCase) Execute() ([]OrderOutputDto, error) {
	orders, err := u.OrderRepository.List()
	if err != nil {
		return []OrderOutputDto{}, err
	}

	var output []OrderOutputDto
	for _, order := range orders {
		output = append(output, OrderOutputDto{
			ID:         order.ID,
			Price:      order.Price,
			Tax:        order.Tax,
			FinalPrice: order.FinalPrice,
		})
	}

	return output, nil
}
