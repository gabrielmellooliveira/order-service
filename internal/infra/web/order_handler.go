package web

import (
	"encoding/json"
	"github.com/gabrielmellooliveira/order-service/internal/entity"
	"github.com/gabrielmellooliveira/order-service/internal/usecase"
	"github.com/gabrielmellooliveira/order-service/pkg/events"
	"net/http"
)

type WebOrderHandler struct {
	OrderRepository entity.OrderRepositoryInterface
	OrderCreated    events.EventInterface
	EventDispatcher events.EventDispatcherInterface
}

func NewWebOrderHandler(
	orderRepository entity.OrderRepositoryInterface,
	orderCreated events.EventInterface,
	eventDispatcher events.EventDispatcherInterface,
) *WebOrderHandler {
	return &WebOrderHandler{
		OrderRepository: orderRepository,
		OrderCreated:    orderCreated,
		EventDispatcher: eventDispatcher,
	}
}

func (h *WebOrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto usecase.OrderInputDto
	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createOrder := usecase.NewCreateOrderUseCase(h.OrderRepository, h.OrderCreated, h.EventDispatcher)
	output, err := createOrder.Execute(dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *WebOrderHandler) List(w http.ResponseWriter, r *http.Request) {
	listOrders := usecase.NewListOrdersUseCase(h.OrderRepository)
	output, err := listOrders.Execute()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
