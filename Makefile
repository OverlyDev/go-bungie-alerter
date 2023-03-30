MAIN := .
OUT := bin/hello

build: format tidy clean
	go build -o $(OUT) $(MAIN)

run: format
	go run $(MAIN)

tidy:
	go mod tidy

format:
	go fmt

clean:
	rm -f bin/*
	rm -rf _context

exec: build
	$(OUT)

minify: clean
	go build -o bin/hello-normal $(MAIN)
	go build -ldflags "-s -w" -o bin/hello-stripped $(MAIN)
	upx -o bin/hello-normal-upx bin/hello-normal
	upx -o bin/hello-stripped-upx bin/hello-stripped
	upx --best --lzma -o bin/hello-normal-upx-bestlzma bin/hello-normal
	upx --best --lzma -o bin/hello-stripped-upx-bestlzma bin/hello-stripped
	ls -lh bin

docker: clean
	mkdir _context
	cp go.* _context/.
	cp *.go _context/.
	cd _context && docker build -t overlydev/go-hello -f ../Dockerfile .