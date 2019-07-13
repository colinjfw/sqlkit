PACKAGES=./db ./encoding ./example

lint:
	gometalinter $(PACKAGES)
.PHONY: lint

test:
	@echo "testing $$SQLKIT_DRIVER $$SQLKIT_CONN"
	GO111MODULE=on go test -v -race -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)
	@echo "--- PASS ---"
.PHONY: test

services:
	docker run --rm -d \
	  --name=sqlkit-mysql \
	  -e MYSQL_USER=travis \
	  -e MYSQL_PASSWORD="" \
	  -e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
	  -e MYSQL_DATABASE=sqlkit \
	  -p 3306:3306 \
	  mysql
	docker run --rm -d \
	  --name=sqlkit-postgres \
	  -e POSTGRES_PASSWORD="" \
	  -e POSTGRES_USER=root \
	  -e POSTGRES_DB=sqlkit \
	  -p 5432:5432 \
	  postgres

.PHONY: services

test-mysql: export SQLKIT_DRIVER = mysql
test-mysql: export SQLKIT_CONN = root@tcp(127.0.0.1:3306)/sqlkit
test-mysql: test
.PHONY: test-mysql

test-postgres: export SQLKIT_DRIVER = postgres
test-postgres: export SQLKIT_CONN = postgres://postgres@127.0.0.1:5432/sqlkit?sslmode=disable
test-postgres: test
.PHONY: test-postgres

update-licenses:
	@hack/update-licenses.sh
.PHONY: update-licenses
