


lint:
	docker run --rm \
		-v $(shell pwd):/go/src/github.com/itimoveev/price-store-test \
		-w /go/src/github.com/itimoveev/price-store-test \
		golangci/golangci-lint:v1.36 golangci-lint run