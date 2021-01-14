build:
	@go build .
run: 
	@go run *.go
help:
	@go run *.go --help
clean:
	rm *.json

run_example:
	@go run *.go --all --dir examples --case snake