.PHONY: all
all: clean build

APP=queue-publisher
APP_EXECUTABLE="./bin/$(APP)"
define DOCKERCOMMAND
docker run -d --name queue-publisher \
-p 8080:8080 \
-e NEW_RELIC_APP_NAME=queue-publisher-dev \
-e NEW_RELIC_LICENSE_KEY=91d33e6993c7986f227461ab3f71b4ecdf3aNRAL \
-e NEW_RELIC_ENABLED=true \
-e AMQP_USERNAME=root \
-e AMQP_PASSWORD=root \
-e AMQP_HOST=some-rabbit \
-e AMQP_PORT=5672 \
-e AMQP_VHOST=test \
-e AMARTHACORE_EXCHANGE=amarthacore \
-e AMARTHACORE_EXCHANGE_TYPE=fanout \
--link some-rabbit:some-rabbit \
queue-publisher:latest
endef

clean:
	rm -f $(APP_EXECUTABLE)
	go mod tidy
	go mod vendor

compile:
	mkdir -p bin/
	go build -o $(APP_EXECUTABLE)

build: clean compile

build-docker:
	docker build -f docker/Dockerfile -t $(APP):latest .

run-docker:
	$(DOCKERCOMMAND)

clean-docker:
	docker stop $(APP) && docker rm -f $(APP)

protoc:
	protoc --go_out=plugins=grpc:. delivery/proto/user.proto
	protoc --go_out=plugins=grpc:. --grpc-gateway_out=grpc_api_configuration=delivery/proto/user.yaml:. delivery/proto/user.proto

lint:
	golangci-lint run --out-format checkstyle > lint.xml

run:
	go run cmd/main.go