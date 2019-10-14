GOROOT := $(GOROOT)
OSJS := GOOS="js"
WASM := GOARCH="wasm"

.PHONY: all

all: 	wasm

wasm:
	cp $(GOROOT)/misc/wasm/wasm_exec.html $(PWD)
	cp $(GOROOT)/misc/wasm/wasm_exec.js $(PWD)
	env $(OSJS) $(WASM) go mod tidy
	env $(OSJS) $(WASM) go build -v

pc:
	go mod tidy
	go build -v
compile:
	go tool compile
clean:
	go clean -x github.com/kdevb0x/quikface
	rm -rf $(PWD)/wasm_exec.js $(PWD)/wasm_exec.html


