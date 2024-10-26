.PHONY: run
run:
	go run cmd/api-server/main.go

# https://github.com/mfridman/tparse is needed
.PHONY: test
test:
	go test ./... -json | tparse -all
