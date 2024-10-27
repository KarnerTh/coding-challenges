.PHONY: run
run:
	go run cmd/api-server/main.go

# https://github.com/mfridman/tparse is needed for pretty output
.PHONY: test
test:
	go test ./... -json | tparse -all

.PHONY: createDevice
createDevice:
	@curl -X POST localhost:8080/api/v0/devices \
	  -H "Content-Type: application/json" \
		-d '{"id": "test1", "algorithm": "RSA"}' 

.PHONY: getDevices
getDevices:
	@curl localhost:8080/api/v0/devices
	 
.PHONY: signTransaction
signTransaction:
	@curl -X POST localhost:8080/api/v0/devices/test1/signatures \
	  -H "Content-Type: application/json" \
		-d '{"data": "importantContent"}' 
