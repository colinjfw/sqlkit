PACKAGES=./db ./encoding

test:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
.PHONY: test

docs:
	godoc2md github.com/coldog/sqlkit/encoding > ./encoding/README.md
	godoc2md github.com/coldog/sqlkit/db > ./db/README.md
.PHONY: docs

update-licenses:
	@hack/update-licenses.sh
.PHONY: update-licenses
