.PHONY: fetch
fetch:
	GoEnv=local GO111MODULE=on go run cmd/fetch/main.go 2>&1

.PHONY: createdb
createdb:
	GoEnv=local GO111MODULE=on go run cmd/createdb/main.go

.PHONY: deletedb
deletedb:
	GoEnv=local GO111MODULE=on go run cmd/deletedb/main.go

.PHONY: createtopics
createtopics:
	GoEnv=local GO111MODULE=on go run cmd/kafka/main.go

.PHONY: analyze
analyze:
	GoEnv=local GO111MODULE=on go run cmd/analyze/main.go

.PHONY: reset
reset: deletedb createdb fetch

.PHONY: watch
watch:
	~/go/bin/reflex -r '\.go' -s -- sh -c "go run ./cmd/fetch/main.go"

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
