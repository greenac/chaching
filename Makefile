.PHONY: run
run:
	GO111MODULE=on go run cmd/main.go 2>&1

.PHONY: watch
watch:
	~/go/bin/reflex -r '\.go' -s -- sh -c "go run ./cmd/main.go"

.PHONY: test
test:
	GO111MODULE=on go test $$(go list ./... | grep -v /mocks) 2>&1

.PHONY: cover
cover:
	GO111MODULE=on go test $$(go list ./... | grep -v /integrations | grep -v /mocks | grep -v /docs | grep -v /setup) -coverprofile tmp/cover.out
	go tool cover -html=tmp/cover.out -o tmp/coverage.html
	open tmp/coverage.html

.PHONY: coverfuncs
coverfuncs:
	GO111MODULE=on go test $$(go list ./... | grep -v /integrations | grep -v /mocks | grep -v /docs | grep -v /setup) -covermode=atomic -coverprofile tmp/cover.out
	go tool cover -func tmp/cover.out -o tmp/function-coverage.out
