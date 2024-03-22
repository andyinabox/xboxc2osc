.PHONY: build
build: clean target/release/xboxc2osc target/release/xboxc2midi

.PHONY: run
run: 
	go run ./cmd/xboxc2osc/main.go

.PHONY: clean
clean: 
	rm -rf target

target/release/xboxc2osc:
	go build -o target/release/xboxc2osc ./cmd/xboxc2osc

target/release/xboxc2midi:
	go build -o target/release/xboxc2midi ./cmd/xboxc2midi