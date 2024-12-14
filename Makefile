# Define variables
GO_CMD=go
MAIN=./cmd/bin/main.go

test:
	go test -v ./... -cover

run:
	$(GO_CMD) run $(MAIN) serve-http

hot:
	@echo " >> Installing gin if not installed"
	@go install github.com/codegangsta/gin@latest
	@gin -i -p 9002 -a 9090 --path cmd/bin --build cmd/bin serve-http

goose-create:
# example : make goose-create name=create_users_table
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
ifndef name
	$(error Usage: make goose-create name=<table_name>)
else
	@goose -dir db/migrations create $(name) sql
endif

goose-up:
# example : make goose-up
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations mysql "root:password@tcp(127.0.0.1:8889)/multifinance?parseTime=true" up

goose-down:
# example : make goose-down
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations mysql "root:password@tcp(127.0.0.1:8889)/multifinance?parseTime=true" down

goose-status:
# example : make goose-status
	@echo " >> Installing goose if not installed"
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@goose -dir db/migrations mysql "root:password@tcp(127.0.0.1:8889)/multifinance?parseTime=true" status

seed:
# make seed total=10 table=roles
	$(GO_CMD) run $(MAIN) seed -total=$(total) -table=$(table)

# Mock generation target
generate-mock:
# example : make generate-mock module=customer source=ports/ports.go destination=service/service_mock_test.go package=service
	@echo " >> Installing mockgen if not installed"
	@go install github.com/golang/mock/mockgen@latest
ifndef module
	$(error Usage: make generate-mock module=<module_name> source=<source_file> destination=<destination_file> package=<package_name>)
endif
ifndef source
	$(error Usage: make generate-mock module=<module_name> source=<source_file> destination=<destination_file> package=<package_name>)
endif
ifndef destination
	$(error Usage: make generate-mock module=<module_name> source=<source_file> destination=<destination_file> package=<package_name>)
endif
ifndef package
	$(error Usage: make generate-mock module=<module_name> source=<source_file> destination=<destination_file> package=<package_name>)
endif
	@mockgen -source=internal/module/$(module)/$(source) \
		-destination=internal/module/$(module)/$(destination) \
		-package=$(package)
	@echo " >> Mock generated for module: $(module)"

mock-all:
# make mock-all
	go generate ./...