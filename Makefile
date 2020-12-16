flatb_go:
	flatc --go ./flatbuff/request.fbs
	flatc --go ./flatbuff/response.fbs

build:
	go build ./cmd/eisenbeton-go/main.go

run:
	go run ./cmd/eisenbeton-go/main.go


