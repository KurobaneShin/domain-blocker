build:
	go build -o domain-blocker .

run:
	go run main.go

install:
	sudo bash -c './domain-blocker &'