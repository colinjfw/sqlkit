PACKAGES=./db ./encoding

test:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -cover $(PACKAGES)

docs:
	godoc2md github.com/coldog/sqlkit/encoding > ./encoding/README.md
	godoc2md github.com/coldog/sqlkit/db > ./db/README.md

ci:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
	bash -c 'bash <(curl -s https://codecov.io/bash) -f coverage.txt'
	rm coverage.txt
