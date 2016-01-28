default: src/renroll/*.go src/main/*.go
	./setpath.sh
	go build -i -o ./renroll src/main/main.go
