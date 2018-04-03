PACKAGES=./db ./encoding

lint:
	gometalinter $(PACKAGES)
.PHONY: lint

test:
	@echo "testing $$SQLKIT_DRIVER $$SQLKIT_CONN"
	go test -i $(PACKAGES)
	go test -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
.PHONY: test

services:
	@hack/services.sh
.PHONY: services

test-mysql: export SQLKIT_DRIVER = mysql
test-mysql: export SQLKIT_CONN = "root@tcp(127.0.0.1:3306)/sqlkit"
test-mysql: test
.PHONY: test-mysql

test-postgres: export SQLKIT_DRIVER = postgres
test-postgres: export SQLKIT_CONN = user=postgres dbname=sqlkit sslmode=disable
test-postgres: test
.PHONY: test-postgres

update-licenses:
	@hack/update-licenses.sh
.PHONY: update-licenses
