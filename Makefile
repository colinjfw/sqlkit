PACKAGES=./db ./encoding

license="// Copyright (C) 2018 Colin Walker \
// \
// This software may be modified and distributed under the terms \
// of the MIT license.  See the LICENSE file for details. \
"

test:
	go test -i $(PACKAGES)
	gometalinter $(PACKAGES)
	go test -coverprofile=coverage.txt -covermode=atomic $(PACKAGES)

docs:
	godoc2md github.com/coldog/sqlkit/encoding > ./encoding/README.md
	godoc2md github.com/coldog/sqlkit/db > ./db/README.md

licenses:
	for file in $(shell find ./db ./encoding); do \
		echo $$file \
		# if [ "$(shell head $$file -n 1 | grep 'Copyright')" == "" ]; then \
			# echo "$(license)\n$(shell cat $$file)" > $$file \
		# fi
	done
