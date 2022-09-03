.PHONY: run
run:
	~/go/bin/reflex -r '\.go' -s -- sh -c "go run ./cmd/main.go"

.PHONY: test
test:
	GO111MODULE=on go test $$(go list ./... | grep -v /mocks) 2>&1
