default: dev

install:
	@rm -f $HOME/go/bin/odas && go generate . > ./internal/commit.txt &&  go install .
build:
	@go generate . > ./internal/commit.txt
	@CGO_ENABLED=0 GOOS=linux go build -a -o ./bin/odas .
dev:
	@go generate . > ./internal/commit.txt
	@go build -o ./bin/odas .