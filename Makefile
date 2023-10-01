create-migration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	migrate -path sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" up

migrate-down:
	migrate -path sql/migrations -database "mysql://root:root@tcp(localhost:3306)/orders" down

.PHONY: create-migration migrate migrate-down