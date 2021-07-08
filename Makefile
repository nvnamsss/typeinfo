build:
	@go build .
setup:
	@cp typeinfo /home/vid-user/go/bin/
run: 
	@go run *.go
help:
	@go run *.go --help
clean:
	rm *.json

run_example:
	@go run *.go --all --dir examples --case snake

run_jf1:
	@go run *.go --format jf1 --all --dir examples/jf --case snake --output infos/jf1

run_jf2:
	@go run *.go --format jf2 --all --dir examples/jf --case snake --output infos/jf2

run_pointer:
	@go run *.go --format jf2 --all --dir examples/pointers --case snake --output infos/pointers