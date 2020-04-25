GOROOT := $(GOROOT)
OSJS := GOOS="js"
WASM := GOARCH="wasm"

.PHONY: all

all: 	wasm

wasm:
	cp $(GOROOT)/misc/wasm/wasm_exec.html $(PWD)/assets/wasm_exec.html
	cp $(GOROOT)/misc/wasm/wasm_exec.js $(PWD)/assets/wasm_exec.js
	env $(OSJS) $(WASM) go mod tidy
	env $(OSJS) $(WASM) go build -v

pc:
	go mod tidy
	go build -v ./...

compile:
	go tool compile -m -+ -dynlink -linkobj quikface.so -o quikface.o src/*.go

clean:
	rm -rf $(PWD)/assets/wasm_exec.js $(PWD)/assets/wasm_exec.html
	rm -rf ./cmd/server/server ./src/quikface

