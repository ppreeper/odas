install:
	@rm -f $HOME/go/bin/odas && go generate . > commit.txt &&  go install .
build:
	@go generate . > commit.txt
	@CGO_ENABLED=0 GOOS=linux go build -a -o bin/odas .
