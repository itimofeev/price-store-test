
lint:
	docker run --rm \
		-v $(shell pwd):/go/src/github.com/itimoveev/price-store-test \
		-w /go/src/github.com/itimoveev/price-store-test \
		golangci/golangci-lint:v1.36 golangci-lint run

run-db:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password --name=db postgres:13

stop-db:
	docker rm -f db

deploy:
	docker-compose -f deployments/docker-compose.yml up -d

stop:
	docker-compose -f deployments/docker-compose.yml kill
	docker-compose -f deployments/docker-compose.yml rm -f
