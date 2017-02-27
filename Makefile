.PHONY: build release clean

build:
	go build -o turtle-proxy

release:
	GOOS=linux GOARCH=amd64 go build -o turtle-proxy-linux-x64
	GOOS=darwin GOARCH=amd64 go build -o turtle-proxy-darwin-x64
	GOOS=windows GOARCH=amd64 go build -o turtle-proxy-windows-x64.exe
	tar czvf turtle-proxy-linux-x64.tar.gz turtle-proxy-linux-x64 README.md LICENSE
	tar czvf turtle-proxy-darwin-x64.tar.gz turtle-proxy-darwin-x64 README.md LICENSE
	tar czvf turtle-proxy-windows-x64.tar.gz turtle-proxy-windows-x64.exe README.md LICENSE

clean:
	go clean
	rm -rf turtle-proxy-*
