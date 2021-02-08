
lint:
	docker run --rm \
		-v $(shell pwd):/go/src/github.com/itimoveev/price-store-test \
		-w /go/src/github.com/itimoveev/price-store-test \
		golangci/golangci-lint:v1.36 golangci-lint run

build-image:
	docker build --force-rm=true \
		-t price-store \
		-f build/Dockerfile .

deploy:
	docker-compose -f deployments/docker-compose.yml up -d

stop:
	docker-compose -f deployments/docker-compose.yml kill
	docker-compose -f deployments/docker-compose.yml rm -f
