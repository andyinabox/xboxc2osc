.PHONY: run
run: 
	go run ./cmd/xboxc2osc/main.go

.PHONY: clean
clean: 
	rm -rf target

.PHONY: build
build: clean target/release/xboxc2osc

target/release/xboxc2osc:
	go build -o target/release/xboxc2osc ./cmd/xboxc2osc