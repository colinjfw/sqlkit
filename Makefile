PACKAGES=./db ./encoding

test:
	@gometalinter $(PACKAGES)
	@go test -cover $(PACKAGES)

docs:
	godoc2md github.com/coldog/sqlkit/encoding > ./encoding/README.md
	godoc2md github.com/coldog/sqlkit/db > ./db/README.md
