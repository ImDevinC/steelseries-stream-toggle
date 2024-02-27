build-windows-x64:
	mkdir -p out/windows-x64
	GOOS=windows GOARCH=amd64 go build -o out/windows-x64/steelseries-stream-toggle.exe ./cmd/sst.go

clean:
	rm -rf out

release: build-windows-x64
	zip -j out/steelseries-stream-toggle-${RELEASE_VERSION}.zip out/windows-x64/steelseries-stream-toggle.exe