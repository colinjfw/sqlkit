PACKAGES=./db ./encoding

test:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
.PHONY: test

update-licenses:
	@hack/update-licenses.sh
.PHONY: update-licenses
