.PHONY: build release test clean

build:
	go build -o turtle-proxy

test:
	go test ./...

release: *.tar.gz
	@read -p "tag?" release_tag; \
	for arch in linux darwin windows; do \
		github-release upload \
			--user lpedrosa \
			--repo turtle-proxy \
			--tag $$release_tag \
			--name "turtle-proxy-$$release_tag-$$arch-x64.tar.gz" \
			--file turtle-proxy-$$arch-x64.tar.gz \
			--replace; \
	done

clean:
	go clean
	rm -rf turtle-proxy-*

*.tar.gz: *-x64 *-x64.exe
	tar czvf turtle-proxy-linux-x64.tar.gz turtle-proxy-linux-x64 README.md LICENSE > /dev/null
	tar czvf turtle-proxy-darwin-x64.tar.gz turtle-proxy-darwin-x64 README.md LICENSE > /dev/null
	tar czvf turtle-proxy-windows-x64.tar.gz turtle-proxy-windows-x64.exe README.md LICENSE > /dev/null

*-x64:
	GOOS=linux GOARCH=amd64 go build -o turtle-proxy-linux-x64
	GOOS=darwin GOARCH=amd64 go build -o turtle-proxy-darwin-x64

*-x64.exe:
	GOOS=windows GOARCH=amd64 go build -o turtle-proxy-windows-x64.exe

