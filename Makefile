default: src/renroll/*.go src/main/*.go
	@GOPATH=${GOPATH}:${HOME}/renroll go build -i -o ./renroll src/main/main.go
