PACKAGES=./db ./encoding

test:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -cover $(PACKAGES)

docs:
	godoc2md github.com/coldog/sqlkit/encoding > ./encoding/README.md
	godoc2md github.com/coldog/sqlkit/db > ./db/README.md

ci:
	@go test -coverprofile=coverage_db.txt -covermode=atomic ./db
	@go test -coverprofile=coverage_encoding.txt -covermode=atomic ./encoding
	echo '' > coverage.txt
	cat coverage_db.txt >> coverage.txt
	cat coverage_encoding.txt >> coverage.txt
	rm coverage_db.txt
	rm coverage_encoding.txt
	bash -c 'bash <(curl -s https://codecov.io/bash)'
	rm coverage.txt
