MAIN := .
OUT := bin/bungie-alerter
LINTER := $(shell go env GOPATH)/bin/golangci-lint

build: format tidy clean
	go generate
	go build -o $(OUT) $(MAIN)

windows: format tidy clean
	go generate
	GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(OUT)-windows-amd64.exe $(MAIN)

run: format
	go run $(MAIN) go

debug: format
	go run $(MAIN) -d -l go

tidy:
	go mod tidy

format:
	go fmt

lint:
	GOGC=off $(LINTER) run

clean:
	rm -f bin/*
	rm -rf _context
	rm -rf embeds

exec: build
	$(OUT)

minify: clean
	go generate
	go build -o $(OUT)-normal $(MAIN)
	go build -ldflags "-s -w" -o $(OUT)-stripped $(MAIN)
	upx --best --lzma -o $(OUT)-normal-upx-bestlzma $(OUT)-normal
	upx --best --lzma -o $(OUT)-stripped-upx-bestlzma $(OUT)-stripped
	ls -lh bin

docker: clean
	mkdir _context
	cp go.* _context/.
	cp *.go _context/.
	cd _context && docker build -t overlydev/go-hello -f ../Dockerfile .