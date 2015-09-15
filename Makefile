default: src/renroll/*.go src/main/*.go
	go build -i -o ./renroll src/main/main.go
