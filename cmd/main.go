package main

import (
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gabrielmellooliveira/order-service/configs"
	"github.com/gabrielmellooliveira/order-service/internal/infra/event/handler"
	"github.com/gabrielmellooliveira/order-service/internal/infra/graph"
	"github.com/gabrielmellooliveira/order-service/internal/infra/grpc/pb"
	"github.com/gabrielmellooliveira/order-service/internal/infra/grpc/service"
	"github.com/gabrielmellooliveira/order-service/internal/infra/web/webserver"
	"github.com/gabrielmellooliveira/order-service/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(config.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	rabbitMQChannel := getRabbitMQChannel(config.RabbitMQUrl)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("order_created", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrdersUseCase(db)

	webServer := webserver.NewWebServer(config.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)

	webServer.AddHandler("/order", webOrderHandler.Create)
	webServer.AddHandler("/orders", webOrderHandler.List)

	go webServer.Start()

	grpcServer := grpc.NewServer()

	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)

	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.GRPCServerPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrdersUseCase:  *listOrdersUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	http.ListenAndServe(":"+config.GraphQLServerPort, nil)
}

func getRabbitMQChannel(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}
