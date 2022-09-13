LIPO = /usr/bin/x86_64-apple-darwin-lipo

GIT_COMMIT_HASH := $(shell git rev-list -1 HEAD)
LDFLAGS_COMMON := -s -w
all: clean \
	grsical-windows-amd64 grsical-linux-amd64 grsical-linux-arm64 grsical-darwin-amd64 grsical-darwin-arm64 \
	grsicalsrv-windows-amd64 grsicalsrv-linux-amd64 grsicalsrv-linux-arm64 grsicalsrv-darwin-amd64 grsicalsrv-darwin-arm64

grsical-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X grs-ical/internal/grsical.version=$(GIT_COMMIT_HASH)" -o build/grsical-windows-amd64.exe grs-ical/cmd/grsical

grsical-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X grs-ical/internal/grsical.version=$(GIT_COMMIT_HASH)" -o build/grsical-linux-amd64 grs-ical/cmd/grsical

grsical-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON) -X grs-ical/internal/grsical.version=$(GIT_COMMIT_HASH)" -o build/grsical-linux-arm64 grs-ical/cmd/grsical

grsical-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON) -X grs-ical/internal/grsical.version=$(GIT_COMMIT_HASH)" -o build/grsical-darwin-amd64 grs-ical/cmd/grsical

grsical-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON) -X grs-ical/internal/grsical.version=$(GIT_COMMIT_HASH)" -o build/grsical-darwin-arm64 grs-ical/cmd/grsical

grsicalsrv-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/grsicalsrv-windows-amd64.exe grs-ical/cmd/grsicalsrv

grsicalsrv-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/grsicalsrv-linux-amd64 grs-ical/cmd/grsicalsrv

grsicalsrv-linux-arm64:
	GOOS=linux GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/grsicalsrv-linux-arm64 grs-ical/cmd/grsicalsrv

grsicalsrv-darwin-amd64:
	GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/grsicalsrv-darwin-amd64 grs-ical/cmd/grsicalsrv

grsicalsrv-darwin-arm64:
	GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS_COMMON)" -o build/grsicalsrv-darwin-arm64 grs-ical/cmd/grsicalsrv

merge-macos-binary:
	$(LIPO) -create build/grsical-darwin-amd64 build/grsical-darwin-arm64 -o build/grsical-darwin-universal
	$(LIPO) -create build/grsicalsrv-darwin-amd64 build/grsicalsrv-darwin-arm64 -o build/grsicalsrv-darwin-universal

docker-all:
	DOCKER_BUILDKIT=1 docker build --file Dockerfile --output build .

clean:
	-rm -f build/*

.PHONY: all clean \
	grsical-windows-amd64 grsical-linux-amd64 grsical-linux-arm64 grsical-darwin-amd64 grsical-darwin-arm64 \
	grsicalsrv-windows-amd64 grsicalsrv-linux-amd64 grsicalsrv-linux-arm64 grsicalsrv-darwin-amd64 grsicalsrv-darwin-arm64 \
	docker-all